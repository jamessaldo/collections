import mysql from 'mysql2/promise';
import database from '@configs/database';
// import env from '@configs/environment';

const { mysql_host, mysql_port, mysql_user, mysql_password, mysql_name } = database;

const pool = mysql.createPool({
    host: mysql_host,
    port: Number(mysql_port),
    user: mysql_user,
    password: mysql_password,
    database: mysql_name,
    // debug: env.NODE_ENV === "production" ? false : true,
});

pool.on('connection', () => {
    console.log('Connected to MySQL DB');
});


const query = async <T extends mysql.RowDataPacket[][] | mysql.RowDataPacket[] | mysql.OkPacket | mysql.OkPacket[] | mysql.ResultSetHeader>
    (sql: string, params: Array<string | number>): Promise<[T, mysql.FieldPacket[]]> => {
    // Get a connection from the pool
    const connection = await pool.getConnection();

    try {
        // Execute the query
        const [results, fields] = await connection.execute(sql, params);

        // Return the results and release the connection back to the pool
        return [results as T, fields];
    } finally {
        // Always release the connection back to the pool when done
        connection.release();
    }
}

export default { query };
