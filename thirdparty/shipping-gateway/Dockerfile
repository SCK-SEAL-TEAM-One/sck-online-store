FROM bbyars/mountebank:2.9.1
COPY ./imposters.json /imposters/imposters.json
COPY ./response.js /imposters/response.js
CMD [ "start", "--configfile", "/imposters/imposters.json", "--allowInjection" ]
