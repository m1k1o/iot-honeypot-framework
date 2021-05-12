#!/bin/bash

cp supervisord.{conf,run}

function generate_entry {
    echo ''
    echo "[program:fwd-${1}-${3}]"
    echo "command=/usr/src/proxy.sh '${1}' '${HOST}' '${2}' '${3}'"
    echo 'autostart=true'
    echo 'autorestart=true'
    echo 'startretries=5'
    echo 'numprocs=1'
    echo 'startsecs=0'
    #echo 'user=%(ENV_USER)s' # Run as root (bind ports <1024)
    echo 'process_name=%(program_name)s_%(process_num)02d'
    echo 'stderr_logfile=/dev/stdout'
    echo 'stderr_logfile_maxbytes=0'
    echo 'stdout_logfile=/dev/stderr'
    echo 'stdout_logfile_maxbytes=0'
} >> /usr/src/supervisord.run

function tcpdump_entry {
    echo ''
    echo "[program:tcpdump]"
    echo "command=/usr/sbin/tcpdump host '${HOST}' -i eth0 -s 65535 -w /var/log/capture/tcpdump -W 48 -G 1800 -C 100 -K -n"
    echo 'autostart=true'
    echo 'autorestart=true'
    echo 'startretries=5'
    echo 'numprocs=1'
    echo 'startsecs=0'
    #echo 'user=%(ENV_USER)s' # Run as root (bind ports <1024)
    echo 'process_name=%(program_name)s_%(process_num)02d'
    echo 'stderr_logfile=/dev/stdout'
    echo 'stderr_logfile_maxbytes=0'
    echo 'stdout_logfile=/dev/stderr'
    echo 'stdout_logfile_maxbytes=0'
} >> /usr/src/supervisord.run

SAVEIFS=${IFS}     # Save current IFS
IFS=$'\n'          # Change IFS to new line
PORTS=(${PORTS})   # Split to array $PORTS
IFS=${SAVEIFS}     # Restore IFS

for (( i=0; i<${#PORTS[@]}; i++ ))
do
    ROW=(${PORTS[$i]})
    PROTOCOL=${ROW[0]}
    TARGET_PORT=${ROW[1]}
    PUBLISHED_PORT=${ROW[2]}

    if [[ ${PROTOCOL:0:1} == "#" ]]
    then
        continue
    fi

    generate_entry "${PROTOCOL}" "${TARGET_PORT}" "${PUBLISHED_PORT}"
done

if [[ ! -z "${CAPTURE}" ]]
then
    tcpdump_entry
fi

#
# start supervisor
/usr/bin/supervisord -c /usr/src/supervisord.run
