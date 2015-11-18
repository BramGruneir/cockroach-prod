#!/bin/bash
# This is a simple wrapper script to download and run sql.test.
# We intentionally do not fail on errors as we want to write the DONE
# file, even on failures.
set -x

LOG_DIR="logs"
BINARY="sql.test"
mkdir -p ${LOG_DIR}

time ./${BINARY} --test.run=TestLogic -d "test/index/*/*/*.test" > \
  ${LOG_DIR}/${BINARY}.STDOUT 2> ${LOG_DIR}/${BINARY}.STDERR < /dev/null

# Find the target of the symlink. It contains the build sha.
binary_name=$(readlink ${BINARY} || echo ${BINARY})

# SECONDS is the time since the shell started. This is a good approximation for now,
# more details are in the output.
echo "time: ${SECONDS}" > ${LOG_DIR}/DONE
echo "binary: ${binary_name}" >> ${LOG_DIR}/DONE
