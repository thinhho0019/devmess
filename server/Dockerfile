# Sử dụng image Postgres chính thức
FROM postgres:15

# Thiết lập biến môi trường mặc định (có thể override khi run container)
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=admin123
ENV POSTGRES_DB=mydb

# Copy các script init nếu cần
# COPY ./init.sql /docker-entrypoint-initdb.d/

# Expose port
EXPOSE 5432
