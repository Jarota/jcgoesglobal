function renderCarousel(carousel, post) {
  let zoomedHtml = ''
  post.pics.forEach((path) => {
    zoomedHtml += `
      <img src="${path}" class="large-pic">
    `
  })
  carousel.innerHTML = zoomedHtml
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
        <p class="date">${new Date(post.created_at).toDateString()}</p>
        <p class="caption">${post.caption}</p>
        ${sig}
      </div>`

  const pics = renderPics(post.id, post.pics)
  if (pics) {
    postDiv.appendChild(pics)
  }

  return postDiv
}

function renderPics(id, paths) {
  if (!paths || paths.length === 0) {
    return null
  }

  let picsDiv = document.createElement('div')
  picsDiv.className = 'pics'
  picsDiv.dataset.zoom = id

  let degrees = -10
  let z = -1
  paths.forEach((path) => {
    let container = document.createElement('div')

    const containerId = `container-${path}`
    container.id = containerId
    container.className = 'pic-cntr'
    container.style.transform = `rotate(${degrees}deg)`
    container.style.zIndex = z

    let loader = document.createElement('div')
    const loaderId = `loader-${path}`
    loader.id = loaderId
    loader.className = 'loader'
    container.appendChild(loader)

    let pic = document.createElement('img')
    pic.className = 'pic'
    pic.dataset.zoom = id
    pic.onload = () => {
      let loader = document.getElementById(loaderId)
      document.getElementById(containerId).replaceChild(pic, loader)
    }
    pic.src = path

    picsDiv.appendChild(container)

    degrees = (degrees > 0 ? degrees + 10 : degrees - 10) * -1
    z -= 1
  })

  return picsDiv
}

export {renderCarousel, hideCarousel, renderTimeline, renderPost}
