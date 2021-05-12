FROM python:3-slim

# https://github.com/darkarnium/kako
WORKDIR /usr/app
COPY . .

RUN pip3 install click requests cerberus setuptools
RUN python setup.py install

ENTRYPOINT [ "/usr/local/bin/kako-master" ]
CMD [ "--configuration-file", "sample/kako.yaml" ]
