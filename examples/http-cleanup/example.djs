let http = require('http')
let fs = require('fs')

function fetchUserData() {
  // Create temporary file for logging
  let logFile = fs.openSync('request.log', 'w')
  defer fs.closeSync(logFile)
  defer console.log('✓ Log file closed')

  // Create HTTP agent
  let agent = new http.Agent({ keepAlive: true, maxSockets: 10 })
  defer agent.destroy()
  defer console.log('✓ HTTP agent destroyed\n')

  fs.writeSync(logFile, `[${new Date().toISOString()}] Starting HTTP request\n`)
  console.log('Starting HTTP cleanup example...')
  console.log('Fetching data from local server (http://localhost:3000)\n')

  // Make synchronous-ish request using callback
  makeRequest('http://localhost:3000/api/users/123', agent, logFile, function(userData, error) {
    if (error) {
      console.error('Failed to fetch user data:', error.message)
      console.error('Make sure to run: npm run server')
      return
    }

    console.log(`✓ User: ${userData.name} (${userData.email})`)

    // Make second request for posts
    makeRequest(`http://localhost:3000/api/users/${userData.id}/posts`, agent, logFile, function(posts, error) {
      if (error) {
        console.error('Failed to fetch posts:', error.message)
        return
      }

      console.log(`✓ Found ${posts.length} posts:\n`)
      posts.forEach(function(post) {
        console.log(`  - ${post.title}`)
      })

      console.log('\n✓ Requests completed')
      console.log('✓ Cleanup will happen via defer statements...\n')
    })
  })
}

// Simplified synchronous-style request function
function makeRequest(url, agent, logFile, callback) {
  let urlObj = require('url').parse(url)
  
  fs.writeSync(logFile, `[${new Date().toISOString()}] GET ${url}\n`)

  let options = {
    hostname: urlObj.hostname,
    port: urlObj.port,
    path: urlObj.path,
    method: 'GET',
    agent: agent,
    timeout: 5000
  }

  let req = http.request(options, function(res) {
    let data = ''

    res.on('data', function(chunk) {
      data = data + chunk
    })

    res.on('end', function() {
      fs.writeSync(logFile, `[${new Date().toISOString()}] Response: ${res.statusCode}\n`)
      
      if (res.statusCode >= 200 && res.statusCode < 300) {
        callback(JSON.parse(data), null)
      } else {
        callback(null, new Error(`HTTP ${res.statusCode}: ${res.statusMessage}`))
      }
    })
  })

  req.on('error', function(err) {
    fs.writeSync(logFile, `[${new Date().toISOString()}] Error: ${err.message}\n`)
    callback(null, err)
  })

  req.on('timeout', function() {
    req.destroy()
    callback(null, new Error('Request timeout'))
  })

  req.end()
}

// Run example
fetchUserData()

// Wait a bit to let async operations complete, then show defer cleanup
setTimeout(function() {
  console.log('Note: defer statements executed when fetchUserData() returned')
  console.log('Check request.log for the full request log')
}, 1000)