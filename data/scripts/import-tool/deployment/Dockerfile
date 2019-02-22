FROM python:3

RUN pip install --upgrade pip && \
    pip install scipy pandas==0.23.4 matplotlib sklearn

COPY ./ root/import-tool

WORKDIR /root/import-tool
CMD [ "python", "./using_clustering.py" ]