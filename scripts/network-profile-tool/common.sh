# This file is meant to be sourced by the scripts step1.sh and step2.sh.

# Prints usage of script.
function usage {
    echo -e "Wrong arguments, usage: bash $0 MANDATORY [OPTIONAL]"
}

# Ensures that the number of passed args are at least equals to the declared number of mandatory args.
# It also handles the special case of the -h or --help arg.
# Arguments of the function:
#   $1: number of arguments passed to the scripts, e.g. "5"
#   $2: first argument passed to the script, e.g. "-h"
#   $3: number of mandatory arguments, e.g. "3"
function margs_precheck {
  if [ "$1" -lt 1 ]; then
    usage
    help
    exit 1
  fi

	if [ "$2" ] && [ "$1" -lt "$3" ]; then
		if [ "$2" == "--help" ] || [ "$2" == "-h" ]; then
			help
			exit
		else
	    usage
	    help
	    exit 1
		fi
	fi
}

# Ensures that all the mandatory args are not empty.
# Arguments of the function:
#   $1: number of mandatory arguments, e.g. "3"
#   $2,3,...: all mandatory arguments
function margs_check {
	if [ $# -lt $(("$1" + 1)) ]; then
	    usage
	    help
	    exit 1 # error
	fi
}

# Exports convenience variables.
# Arguments of the function:
#   $1: network name
#   $2: node index
#   $3: medco version override (can be empty)
function export_variables {
  source ../../scripts/versions.sh "$3"

  PROFILE_NAME="network-$1-node$2"
  SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
  COMPOSE_FOLDER="${SCRIPT_FOLDER}/../../deployments/${PROFILE_NAME}"
  CONF_FOLDER="${COMPOSE_FOLDER}/configuration"
  MEDCO_BIN=(docker run -v "$CONF_FOLDER:/medco-configuration" -u "$(id -u):$(id -g)" "${MEDCO_DOCKER_TAG}")

  export PROFILE_NAME SCRIPT_FOLDER COMPOSE_FOLDER CONF_FOLDER MEDCO_BIN
}

# Performs check of necessary dependencies.
function dependency_check {
  echo "### Check of dependencies, script will abort if not found"
  which docker openssl
}

# Check validity of network name
function check_network_name {
  if [[ ! $1 =~ ^[a-zA-Z0-9-]+$ ]]; then
      echo "Network name must only contain basic characters (a-z, A-Z, 0-9, -)"
      exit 1
  fi
}
