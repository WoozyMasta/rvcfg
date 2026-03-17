#define MAKE(A, B) A + B
class CfgTest
{
    value1 = MAKE(1);
    value2 = MAKE(1, 2, 3);
    value3 = MAKE(1,2);
    value4 = MAKE(1/*x,y*/,2);
};
