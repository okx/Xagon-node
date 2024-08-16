// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.4;

contract Sha {
    function hash() public {
        sha256("hello world");
    }
}
