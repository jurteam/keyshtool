#!/bin/bash

set -eo pipefail

./keyshtool -parts 5 -threshold 3 -f testdata/secret_00 -output ./TESTS/TEST1 split && echo OK || echo KO
./keyshtool -parts 2 -threshold 1 -f testdata/secret_00 -output ./TESTS/TEST2 split 2>/dev/null && echo KO || echo OK
./keyshtool -parts 3 -threshold 3 -f testdata/secret_01 -output ./TESTS/TEST3 split && echo OK || echo KO

./keyshtool combine ./TESTS/TEST1/PARTS/* | cut -d ':' -f2 | sed 's/^ //' >check/testresult_00 && echo OK || echo KO
./keyshtool combine ./TESTS/TEST3/PARTS/* | cut -d ':' -f2 | sed 's/^ //' >check/testresult_01 && echo OK || echo KO

cmp=$(diff -u testdata/secret_00 check/testresult_00)
[ $? -eq 0 ] && [ -z $cmp ] && echo OK || echo KO

cmp=$(diff -u testdata/secret_01 check/testresult_01)
[ $? -eq 0 ] && [ -z $cmp ] && echo OK || echo KO
