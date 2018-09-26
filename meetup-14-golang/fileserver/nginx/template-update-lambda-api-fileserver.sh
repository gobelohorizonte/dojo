#!/bin/bash

#
# autor: @jeffotoni
# about: Script to deploy our applications lambda-api-server
# date:  21/09/2018
# since: Version 0.1
#
DOCKER_NAME='fileserver'

# nome da imagem
DOCKER_IMAGE="jeffotoni/$DOCKER_NAME"

# local beta
#DOCKER_NETWORK='--net s31'
SUDO=""

DB_HOST=$(cat .dbhost)

# producao
DOCKER_NETWORK=''

# volume onde ira econtrar o projeto
# beta local
#DOCKER_VOLUME='/PATH'

# producao
DOCKER_VOLUME='/PATH_VOLUME'

# id container blue
CONTAINER_ID_BLUE=$($SUDO docker ps -q -f name=blue)

# id container green
CONTAINER_ID_GREEN=$($SUDO docker ps -q -f name=green)

PORTA_BLUE="5000"

PORTA_GREEN="5001"

# script php a ser executado
APP_PHP='PATH_APP'

PATH_BLUE_GREEN='/etc/nginx/fileserver-bluegreen.conf'

# buscando status atual do nosso blue green
# set $activeBackend blue;
BLUE_GREEN=$(cat ${PATH_BLUE_GREEN} | grep blue)

echo ""

if [ ! -z "$BLUE_GREEN" ]; then

	echo ""
	echo "########## Esta online é o Blue ##########"
	echo "Deploy no Green [${CONTAINER_ID_GREEN}]"


	# aguarde
	sleep 1

	# derrubar o docker green
	# fazer pull da imagem para atualizar
	# subir docker green
	# e finish
	docker stop ${CONTAINER_ID_GREEN}
	docker rm ${CONTAINER_ID_GREEN}
	docker pull ${DOCKER_IMAGE}

	docker run --net s31 -d -e PORT_SERVER="${PORTA_GREEN}" -e DB_HOST_1="${DB_HOST}" -e AWS_REGION="us-east-1" -e TZ="America/Sao_Paulo" -p ${PORTA_GREEN}:${PORTA_GREEN} -v /golang/storage:/storage  --restart=always --name fileserver-${PORTA_GREEN}-green ${DOCKER_IMAGE}

	# aguarde
	sleep 1

	# alterar para green
	# reload no nginx reload
	echo "set \$activeBackend fileservergreen;" > $PATH_BLUE_GREEN
	echo "Alterado para green"

	# reload no nginx -s reload para atualizar as configuracoes
	/etc/init.d/nginx reload

else 

	echo ""
	echo "########## Esta online é o Green ##########"
	echo "Deploy no Blue [${CONTAINER_ID_BLUE}]"

	sleep 1

	docker stop ${CONTAINER_ID_BLUE}
	docker rm ${CONTAINER_ID_BLUE}
	docker pull ${DOCKER_IMAGE}

	docker run --net s31 -d -e PORT_SERVER="${PORTA_BLUE}" -e DB_HOST_1="${DB_HOST}" -e AWS_REGION="us-east-1" -e TZ="America/Sao_Paulo" -p ${PORTA_BLUE}:${PORTA_BLUE} -v /golang/storage:/storage --restart=always --name fileserver-${PORTA_BLUE}-blue ${DOCKER_IMAGE}

	# alterar para green
	echo "set \$activeBackend fileserverblue;" > $PATH_BLUE_GREEN
	echo "Alterado para blue"

	# reload no nginx -s reload para atualizar as configuracoes
	/etc/init.d/nginx reload
fi

exit 0

