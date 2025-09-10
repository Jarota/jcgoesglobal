import {renderCarousel, hideCarousel, renderPost} from './render.js'

let captionInput = document.getElementById('caption')
let fileInput = document.getElementById('image-upload')
let preview = document.getElementById('preview')
let carousel = document.getElementById('carousel')

// obj to preview
let post = {
  id:        'preview-id',
  author:    '',
  caption:   captionInput.value,
  pics:      [],
  createdAt: ''
}

captionInput.addEventListener('change', updatePreviewCaption)
function updatePreviewCaption() {
  post.caption = captionInput.value
  preview.innerHTML = renderPost(post)
}

fileInput.addEventListener('change', updatePreviewPics)
function updatePreviewPics() {
  post.pics = []
  for (const file of fileInput.files) {
    const reader = new FileReader()

    reader.onload = function (e) {
      post.pics.push(e.target.result)
      preview.innerHTML = renderPost(post)
    }
    reader.readAsDataURL(file)
  }
}

document.addEventListener('click', (e) => {
  // Show the carousel when clicking on pics
  const postToZoom = e.target.dataset.zoom
  if (postToZoom) {
    renderCarousel(carousel, post)
    return
  }

  // Hide the carousel when 'clicking away'
  if (e.target.id !== 'carousel' && !e.target.classList.contains('large-pic')) {
    hideCarousel()
    return
  }
})
