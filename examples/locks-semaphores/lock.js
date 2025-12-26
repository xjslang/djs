// Simple lock implementation for demonstration purposes
class Lock {
  constructor(name) {
    this.name = name
    this.locked = false
    this.queue = []
  }

  async acquire() {
    if (!this.locked) {
      this.locked = true
      console.log(`🔒 Lock "${this.name}" acquired`)
      return
    }

    // Wait in queue
    await new Promise(resolve => {
      this.queue.push(resolve)
    })
    this.locked = true
    console.log(`🔒 Lock "${this.name}" acquired (from queue)`)
  }

  release() {
    console.log(`🔓 Lock "${this.name}" released`)
    if (this.queue.length > 0) {
      let resolve = this.queue.shift()
      resolve()
    } else {
      this.locked = false
    }
  }
}

// Simple semaphore implementation
class Semaphore {
  constructor(name, maxConcurrent) {
    this.name = name
    this.maxConcurrent = maxConcurrent
    this.current = 0
    this.queue = []
  }

  async acquire() {
    if (this.current < this.maxConcurrent) {
      this.current++
      console.log(`🎫 Semaphore "${this.name}" acquired (${this.current}/${this.maxConcurrent})`)
      return
    }

    // Wait in queue
    await new Promise(resolve => {
      this.queue.push(resolve)
    })
    this.current++
    console.log(`🎫 Semaphore "${this.name}" acquired from queue (${this.current}/${this.maxConcurrent})`)
  }

  release() {
    this.current--
    console.log(`🎟️  Semaphore "${this.name}" released (${this.current}/${this.maxConcurrent})`)
    
    if (this.queue.length > 0) {
      let resolve = this.queue.shift()
      resolve()
    }
  }
}

// Lock manager
let locks = new Map()
let semaphores = new Map()

function acquireLock(resourceId) {
  if (!locks.has(resourceId)) {
    locks.set(resourceId, new Lock(resourceId))
  }
  return locks.get(resourceId).acquire().then(function() {
    return locks.get(resourceId)
  })
}

function acquireSemaphore(name, maxConcurrent) {
  if (!semaphores.has(name)) {
    semaphores.set(name, new Semaphore(name, maxConcurrent))
  }
  return semaphores.get(name).acquire().then(function() {
    return semaphores.get(name)
  })
}

module.exports = { acquireLock, acquireSemaphore }
