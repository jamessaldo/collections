import { Pool, QueryResult } from 'pg';
import database from '@configs/database';

const { pg_driver, pg_host, pg_port, pg_user, pg_password, pg_name } = database;

const pool = new Pool({
    connectionString: `${pg_driver}://${pg_user}:${pg_password}@${pg_host}:${pg_port}/${pg_name}`,
});

pool.on('connect', () => {
    console.log('Connected to Postgres DB');
});

const query = async (sql: string, params: Array<string | number>): Promise<QueryResult> => {
    // return await pool.query(sql, params);
    const client = await pool.connect();

    try {
        // Execute the query
        return await client.query(sql, params);
    } finally {
        // Always release the client back to the pool when done
        client.release();
    }
};

export default { query };
