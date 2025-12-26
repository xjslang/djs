let http = require('http')
let fs = require('fs')

// Global resources

async function fetchUserData() {
  let logFile = fs.openSync('request.log', 'w')
  let agent = new http.Agent({ keepAlive: true, maxSockets: 10 })
  defer {
    console.log('\nCleaning up resources...')
    agent.destroy()
    console.log('✓ HTTP agent destroyed')
    fs.closeSync(logFile)
    console.log('✓ Log file closed\n')
  }

  fs.writeSync(logFile, `[${new Date().toISOString()}] Starting HTTP requests\n`)
  console.log('Starting HTTP cleanup example...')
  console.log('Fetching data from local server (http://localhost:3000)\n')

  // Fetch user data with await and or
  let userResponse = await makeRequest('http://localhost:3000/api/users/123', agent, logFile) or |err| {
    console.error('Failed to fetch user data:', err.message)
    console.error('Make sure to run: npm run server')
    return
  }

  let userData = JSON.parse(userResponse.data)
  console.log(`✓ User fetched: ${userData.name} (${userData.email})`)

  // Fetch user's posts with await and or
  let postsResponse = await makeRequest(`http://localhost:3000/api/users/${userData.id}/posts`, agent, logFile) or |err| {
    console.error('Failed to fetch posts:', err.message)
    return
  }

  let posts = JSON.parse(postsResponse.data)
  console.log(`✓ Found ${posts.length} posts:\n`)

  // Display each post
  posts.forEach(function(post) {
    console.log(`  - ${post.title}`)
  })

  console.log('\n✓ All requests completed successfully!')
}

// Helper function to make HTTP requests (returns Promise)
function makeRequest(url, agent, logFile) {
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

  return new Promise(function(resolve, reject) {
    let req = http.request(options, function(res) {
      let data = ''

      res.on('data', function(chunk) {
        data = data + chunk
      })

      res.on('end', function() {
        fs.writeSync(logFile, `[${new Date().toISOString()}] Response: ${res.statusCode}\n`)
        
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve({ data: data, statusCode: res.statusCode })
        } else {
          reject(new Error(`HTTP ${res.statusCode}: ${res.statusMessage}`))
        }
      })
    })

    req.on('error', function(err) {
      fs.writeSync(logFile, `[${new Date().toISOString()}] Error: ${err.message}\n`)
      reject(err)
    })

    req.on('timeout', function() {
      req.destroy()
      reject(new Error('Request timeout'))
    })

    req.end()
  })
}

// Run the async example
(async function() {
  fetchUserData()
  console.log('Note: Check request.log for the full request log')
})()