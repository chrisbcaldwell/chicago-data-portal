const Pool = require('pg').Pool
const pool = new Pool({
  user: 'Postgres',
  host: 'localhost',
  database: 'chicago',
  password: 'italianbeef',
  port: 5434,
})

