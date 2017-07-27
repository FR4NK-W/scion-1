# Copyright 2017 ETH Zurich
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""
:mod:`extn` --- SCION Origin Validation extension
=================================================================
"""
# Stdlib

# SCION
from lib.crypto.symcrypto import mac
from lib.packet.opt.defines import (
    OPTLengths,
    OPTValidationError, OPTMode)
from lib.packet.opt.base_ext import SCIONOriginPathTraceBaseExtn
from lib.util import hex_str, Raw


class SCIONOriginValidationExtn(SCIONOriginPathTraceBaseExtn):
    """
    Implementation of SCION Origin Validation extension.

    OV extension Header

    0B       1        2        3        4        5        6        7
    +--------+--------+--------+--------+--------+--------+--------+--------+
    | xxxxxxxxxxxxxxxxxxxxxxxx |  Meta  |            Timestamp              |
    +--------+--------+--------+--------+--------+--------+--------+--------+
    |                               DataHash...                             |
    +--------+--------+--------+--------+--------+--------+--------+--------+
    |                            ...DataHash                                |
    +--------+--------+--------+--------+--------+--------+--------+--------+
    |                               Session ID...                           |
    +--------+--------+--------+--------+--------+--------+--------+--------+
    |                            ...Session ID                              |
    +--------+--------+--------+--------+--------+--------+--------+--------+
    |                                  OV_i ...                             |
    +--------+--------+--------+--------+--------+--------+--------+--------+
    |                               ...OV_i  (var length)                   |
    +--------+--------+--------+--------+--------+--------+--------+--------+

    """
    NAME = "SCIONOriginValidationExtn"

    def __init__(self, raw=None):  # pragma: no cover
        """
        :param bytes raw: Raw data holding DataHash, SessionID and OVs
        """
        self.OVs = []
        super().__init__(raw)

    def _parse(self, raw):
        """
        Parse payload to extract values.
        :param bytes raw: raw payload.
        """
        data = Raw(raw, self.NAME)
        super()._parse(data)

        self.meta = bytes([data.pop(OPTLengths.META)])
        self.mode = int(self.meta[0] >> 6)
        self.path_index = int(self.meta[0] & 0x3f)
        self.timestamp = data.pop(OPTLengths.TIMESTAMP)
        self.datahash = data.pop(OPTLengths.DATAHASH)
        self.sessionID = data.pop(OPTLengths.SESSIONID)
        all_ovs = data.pop()
        self.OVs = []
        if len(all_ovs) % 16 == 0:  # check we got valid OVs
            self.ov_count = len(all_ovs)//16
            for ov_index in range(len(all_ovs)):
                self.OVs.append(bytes(all_ovs[ov_index*OPTLengths.OVs:(ov_index+1)*OPTLengths.OVs]))

    @classmethod
    def ov_from_values(cls, ov_list):
        return b''.join(ov_list)

    @classmethod
    def from_values(cls, mode, path_index, timestamp, datahash, sessionID, OVs):  # pragma: no cover
        """
        Construct extension.

        :param int mode: The mode of the extension, 0 <= mode <= 3
        :param int path_index: The path_index for the OV field
        :param bytes timestamp: The timestamp when the extension was created
        :param bytes datahash: The hash of the payload
        :param bytes sessionID: The session ID of the extension.
        :param []bytes OVs: The Origin Validation Fields for the extension
        :returns: The created instance.
        :rtype: SCIONOriginValidationExtn
        :raises: None
        """
        inst = cls()
        init_length = OPTLengths.TOTAL[OPTMode.ORIGIN_VALIDATION_ONLY]+16*len(OVs)
        inst._init_size(inst.bytes_to_hdr_len(init_length)-1)
        assert 0 <= mode <= 3
        inst.mode = mode
        inst.path_index = path_index
        inst.meta = (mode << 6) | path_index
        inst.timestamp = timestamp
        inst.datahash = datahash
        inst.sessionID = sessionID
        inst.OVs = OVs
        return inst

    def pack(self):
        """
        Pack extension into byte string.

        :returns: packed extension.
        :rtype: bytes
        """
        meta = (self.mode << 6) | self.path_index
        packed = [bytes([meta]), self.timestamp, self.datahash, self.sessionID,
                  b"".join(self.OVs)]
        raw = b"".join(packed)
        self._check_len(raw)
        return raw

    def create_ovs_from_path(self, intermediate_key_list, dst_key):
        ov_list = []
        for key in intermediate_key_list:
            ov_list.append(mac(key.drkey, self.datahash))
        ov_list.append(mac(dst_key.drkey, self.datahash))
        return ov_list

    @classmethod
    def check_validity(cls, datahash, sessionID, OVs):
        """
        Check if parameters are valid.

        :param bytes datahash: The hash of the payload
        :param bytes sessionID: The session ID of the extension.
        :param bytes OVs: The Origin Validation Field for the extension
        :raises: OPTValidationError
        """

        if len(datahash) != OPTLengths.DATAHASH:
            raise OPTValidationError("Invalid datahash length %sB. Expected %sB" % (
                len(datahash), OPTLengths.DATAHASH))
        if len(sessionID) != OPTLengths.SESSIONID:
            raise OPTValidationError("Invalid sessionID length %sB. Expected %sB" % (
                len(sessionID), OPTLengths.SESSIONID))
        if len(OVs) % OPTLengths.OVs != 0:
            raise OPTValidationError("Invalid OVs length %sB. Expected a multiple of %sB" % (
                len(OVs), OPTLengths.OVs))

    def __str__(self):
        return "%s(%sB):\n\tDatahash: %s\n\tSessionID: %s\n\tOVs: %s" % (
            self.NAME, len(self), hex_str(self.datahash),
            hex_str(self.sessionID), hex_str(b"".join(self.OVs)))
