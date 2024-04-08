// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract Log0 {
    // opcode 0xa0
    function opLog0() public payable {
        assembly {
            log0(0, 32)
        }
    }

    function opLog00() public payable {
        assembly {
            log0(0, 0)
        }
    }

     function opLog01() public payable {
        assembly {
            log0(0, 28)
        }
    }
}