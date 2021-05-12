#!/bin/sh
cd /usr/share/nginx/html/;

cp install.sh.template install.sh;
sed -i "s/10.8.0.1/$BIND_ADDRESS/g" install.sh;
sed -i "s/7999/$BIND_PORT/g" install.sh;

echo "
You can now run this command on your worker device:

  curl http://$BIND_ADDRESS:$BIND_PORT/install.sh | sudo bash
"
