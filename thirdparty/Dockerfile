FROM bbyars/mountebank:2.9.1
COPY ./imposters.json /imposters/imposters.json
COPY ./bank-gateway /imposters/bank-gateway
COPY ./shipping-gateway /imposters/shipping-gateway
CMD [ "start", "--configfile", "/imposters/imposters.json", "--allowInjection" ]
