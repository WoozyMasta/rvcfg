class CfgMacroMatrix
{
  class Root
  {
#define STRINGIFY(S) #S
#define GLUE(A,B) A##B

    // Expected by DayZ CfgConvert: both lines are removed after preprocessing.
    stringify_result = STRINGIFY(ABC);
    glue_result = GLUE(12,34);

    // Expected: __LINE__ expands to concrete integer line number.
    line_result = __LINE__;

    // Expected: these macros are not expanded in DayZ CfgConvert and stay as literal strings.
    counter_a = __COUNTER__;
    counter_b = __COUNTER__;
    date_str = __DATE_STR__;
    date_iso = __DATE_STR_ISO8601__;
    time_local = __TIME__;
    time_utc = __TIME_UTC__;
    timestamp_utc = __TIMESTAMP_UTC__;
    rand_i32 = __RAND_INT32__;
    rand_u32 = __RAND_UINT32__;
  };
};
