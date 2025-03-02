# Resources:
# ----------
# [1] Initial image: Red Hat Universal Base Image 9 Micro
#     https://catalog.redhat.com/software/containers/ubi9/ubi-micro/615bdf943f6014fa45ae1b58
FROM registry.redhat.io/ubi9/ubi-micro


# 01 - Copy the binary to the image..
COPY ./build/links /usr/local/bin/links

# Copy links.csv file to the container..
RUN mkdir -p /var/lib/links
COPY ./resources/links.csv /var/lib/links/links.csv

# Copy certificates to the container..
RUN mkdir -p /var/certificates
COPY ./resources/certificates/cert.pem ./resources/certificates/privkey.pem ./resources/certificates/chain.pem ./resources/certificates/fullchain.pem /var/certificates

# Check if the binary is executable..
RUN chmod +x /usr/local/bin/links

# Vystavíme port z proměnné prostředí
EXPOSE 8080

# Spustíme binárku
CMD ["/usr/local/bin/links", "--datafile=/var/lib/links/links.csv", "--port=8080", "--cert=/var/certificates/cert.pem", "--key=/var/certificates/privkey.pem"]
