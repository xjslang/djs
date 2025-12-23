// setup.js - Creates a sample SQLite database with users table
let sqlite = require('better-sqlite3');
let db = sqlite('mydata.db');

// Create users table
db.exec(`
  CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  )
`);

// Clear existing data
db.exec('DELETE FROM users');

// Insert sample data
let insert = db.prepare('INSERT INTO users (name, email, active) VALUES (?, ?, ?)');

let users = [
  ['Alice Johnson', 'alice@example.com', 1],
  ['Bob Smith', 'bob@example.com', 1],
  ['Charlie Brown', 'charlie@example.com', 0],
  ['Diana Prince', 'diana@example.com', 1],
  ['Eve Adams', 'eve@example.com', 0],
  ['Frank Miller', 'frank@example.com', 1]
];

for (let i = 0; i < users.length; i++) {
  let user = users[i]
  insert.run(user[0], user[1], user[2]);
}

db.close();

console.log('✓ Database created successfully!');
console.log('✓ Sample data inserted (6 users, 4 active)');
