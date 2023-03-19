const database = {
    pg_driver: process.env.PG_DB_DRIVER || 'postgres',
    pg_host: process.env.PG_DB_HOST || 'localhost',
    pg_port: process.env.PG_DB_PORT || '5432',
    pg_user: process.env.PG_DB_USER || 'postgres',
    pg_password: process.env.PG_DB_PASSWORD || 'postgres',
    pg_name: process.env.PG_DB_NAME || 'postgres',

    mysql_driver: process.env.MYSQL_DB_DRIVER || 'mysql',
    mysql_host: process.env.MYSQL_DB_HOST || 'localhost',
    mysql_port: process.env.MYSQL_DB_PORT || '3306',
    mysql_user: process.env.MYSQL_DB_USER || 'mysql',
    mysql_password: process.env.MYSQL_DB_PASSWORD || 'mysql',
    mysql_name: process.env.MYSQL_DB_NAME || 'mysql',
};

export default database;
