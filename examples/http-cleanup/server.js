// server.js - Simple HTTP server for testing defer and or constructs
let http = require('http');

// Sample data
let users = {
  '123': {
    id: 123,
    name: 'Alice Johnson',
    email: 'alice@example.com',
    active: true
  },
  '456': {
    id: 456,
    name: 'Bob Smith',
    email: 'bob@example.com',
    active: true
  }
};

let posts = {
  '123': [
    { id: 1, userId: 123, title: 'Getting Started with DJS', content: 'Learn about defer and or constructs...' },
    { id: 2, userId: 123, title: 'HTTP Resource Management', content: 'Best practices for HTTP cleanup...' },
    { id: 3, userId: 123, title: 'Error Handling Made Easy', content: 'Using or blocks for elegant error handling...' }
  ],
  '456': [
    { id: 4, userId: 456, title: 'Advanced DJS Patterns', content: 'Deep dive into DJS features...' }
  ]
};

let server = http.createServer(function(req, res) {
  console.log(`${req.method} ${req.url}`);

  // Set CORS headers
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Content-Type', 'application/json');

  // Route: GET /api/users/:id
  let userMatch = req.url.match(/^\/api\/users\/(\d+)$/);
  if (userMatch && req.method === 'GET') {
    let userId = userMatch[1];
    let user = users[userId];
    
    if (user) {
      res.writeHead(200);
      res.end(JSON.stringify(user));
    } else {
      res.writeHead(404);
      res.end(JSON.stringify({ error: 'User not found' }));
    }
    return;
  }

  // Route: GET /api/users/:id/posts
  let postsMatch = req.url.match(/^\/api\/users\/(\d+)\/posts$/);
  if (postsMatch && req.method === 'GET') {
    let userId = postsMatch[1];
    let userPosts = posts[userId] || [];
    
    res.writeHead(200);
    res.end(JSON.stringify(userPosts));
    return;
  }

  // Route: GET /api/health
  if (req.url === '/api/health' && req.method === 'GET') {
    res.writeHead(200);
    res.end(JSON.stringify({ status: 'ok', timestamp: new Date().toISOString() }));
    return;
  }

  // 404 for unknown routes
  res.writeHead(404);
  res.end(JSON.stringify({ error: 'Not found' }));
});

let PORT = 3000;
server.listen(PORT, function() {
  console.log(`✓ HTTP server running on http://localhost:${PORT}`);
  console.log('  Available endpoints:');
  console.log(`  - GET http://localhost:${PORT}/api/users/:id`);
  console.log(`  - GET http://localhost:${PORT}/api/users/:id/posts`);
  console.log(`  - GET http://localhost:${PORT}/api/health`);
  console.log('\nPress Ctrl+C to stop the server');
});

// Graceful shutdown
process.on('SIGINT', function() {
  console.log('\n\nShutting down server...');
  server.close(function() {
    console.log('✓ Server closed');
    process.exit(0);
  });
});
