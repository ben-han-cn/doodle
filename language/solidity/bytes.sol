pragma solidity ^0.4.10;
contract C {
    bytes s;
    function C() {
        s.length = 32 * 4;
        s[31] = 0x1;
        s[63] = 0x2;
        s[95] = 0x3;
        s[127] = 0x4;
    }
}
