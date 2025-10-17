const vidTypes = ['.mp4', '.mov']

function renderCarousel(carousel, post) {
  let zoomedHtml = ''
  post.pics.forEach((pic) => {
    const path = pic.path
    let isVid = vidTypes.some(type => path.includes(type))
    zoomedHtml += isVid ? `
      <video controls src="${path}" class="large-pic">
        video not supported :/
      </video>      
    ` : `
      <img src="${pic.thumbnail}" class="large-pic">
    `
  })
  carousel.innerHTML = zoomedHtml
  
  // reset the scroll for the new carousel
  carousel.scrollTop = 0
  carousel.scrollLeft = 0

  carousel.classList.remove('hidden')
  carousel.classList.add('shown')

  // Also prevent scrolling in the background
  document.body.style.overflow = 'hidden'
  document.documentElement.style.overflow = 'hidden'
}

function hideCarousel() {
  carousel.classList.remove('shown')
  carousel.classList.add('hidden')
  document.body.style.overflow = '' // Re-enable scrolling too
  document.documentElement.style.overflow = ''
}

function renderTimeline(posts) {
  let timeline = document.getElementById('timeline')
  if (posts.length === 0) {
    timeline.innerHTML = `<p>No posts (yet)!</p>`
    return
  }

  let firstPost = true
  posts.forEach((post) => {
    let verticalLine = '<div class="vr"></div>'
    if (firstPost) {
      verticalLine = ''
      firstPost = false
    }

    timeline.innerHTML += `${verticalLine}`
    timeline.appendChild(renderPost(post))
  })
}

function renderPost(post) {
  let postDiv = document.createElement('div')
  postDiv.className = 'post'
  postDiv.id = post.id

  let sig = post.author ? `
    <p class="signature"> - ${post.author}</p>
  ` : ''

  postDiv.innerHTML = 
      `<div class="caption-cntr">
        <p class="date">${new Date(post.date).toDateString()}</p>
        <p class="caption">${post.caption}</p>
        ${sig}
      </div>`

  const pics = renderPics(post.id, post.pics)
  if (pics) {
    postDiv.appendChild(pics)
  }

  return postDiv
}

function renderPics(id, pics) {
  if (!pics || pics.length === 0) {
    return null
  }

  let picsDiv = document.createElement('div')
  picsDiv.className = 'pics'
  picsDiv.dataset.zoom = id

  let degrees = -10
  let z = -1
  pics.forEach((pic) => {
    let container = document.createElement('div')

    const containerId = `container-${pic.id}`
    container.id = containerId
    container.className = 'pic-cntr'
    container.style.transform = `rotate(${degrees}deg)`
    container.style.zIndex = z

    let loader = document.createElement('div')
    const loaderId = `loader-${pic.id}`
    loader.id = loaderId
    loader.className = 'loader'
    container.appendChild(loader)

    const path = pic.path
    const isVid = vidTypes.some(type => path.includes(type))
    if (isVid) {
      renderVid(id, path, containerId, loaderId)
    } else {
      renderPic(id, pic.thumbnail, containerId, loaderId)
    }

    picsDiv.appendChild(container)

    degrees = (degrees > 0 ? degrees + 10 : degrees - 10) * -1
    z -= 1
  })

  return picsDiv
}

function renderVid(postId, path, containerId, loaderId) {
  let vid = document.createElement('video')
  vid.innerHTML = `
    Video not supported :/
  ` // ^ in case embedding doesn't work

  vid.className = 'pic'
  vid.dataset.zoom = postId
  vid.onloadeddata = () => {
    let loader = document.getElementById(loaderId)
    document.getElementById(containerId).replaceChild(vid, loader)
  }
  vid.src = path
}

function renderPic(postId, path, containerId, loaderId) {
  let pic = document.createElement('img')
  pic.className = 'pic'
  pic.dataset.zoom = postId
  pic.onload = () => {
    let loader = document.getElementById(loaderId)
    document.getElementById(containerId).replaceChild(pic, loader)
  }
  pic.src = path
}

export {renderCarousel, hideCarousel, renderTimeline, renderPost}
