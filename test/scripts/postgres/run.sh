#!/bin/bash

set -e

BASEDIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
SCRIPT_FILES="${BASEDIR}/*.sql"
DBNAME=test_db
DBUSER=test_user

main(){
    for origin_script in ${SCRIPT_FILES}; do
        echo "Executing ${origin_script}..."

        script_file_path=$(mktemp)
        script_file_name=$(basename "${script_file_path}")
        script_contents=$(eval "echo \"$(cat ${origin_script})\"")

        echo "${script_contents}" > "${script_file_path}"

        docker cp "${script_file_path}" x1-state-db:"${script_file_path}"
        docker exec x1-state-db bash -c "chmod a+x ${script_file_path} && psql ${DBNAME} ${DBUSER} -v ON_ERROR_STOP=ON --single-transaction -f ${script_file_path}"

        echo "Done"
    done
}

main "${@}"
