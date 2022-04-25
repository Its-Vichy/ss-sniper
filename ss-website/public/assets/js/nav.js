/*  
 * Easy selector helper function
 */
const select = (el, all = false) => {
    el = el.trim()
    if (all) {
        return [...document.querySelectorAll(el)]
    } else {
        return document.querySelector(el)
    }
}

/*
 * Easy event listener function
 */
const on = (type, el, listener, all = false) => {
    let selectEl = select(el, all)
    if (selectEl) {
        if (all) {
            selectEl.forEach(e => e.addEventListener(type, listener))
        } else {
            selectEl.addEventListener(type, listener)
        }
    }
}

/*
 * Easy on scroll event listener 
 */
const onscroll = (el, listener) => {
    el.addEventListener('scroll', listener)
}

let navbarlinks = select('.nav_link', true)
navbarlinks.forEach(navbarlink => {
    navbarlink.onclick = function() {
        select('.active', true).forEach(navlink => {
            navlink.classList.remove('active')
            navlink.classList.add('notactive')
        })
        navbarlink.classList.remove('notactive')
        navbarlink.classList.add('active')
    }
})

// Nav Bar Hambrger

if (window.innerWidth < 425) {
    document.getElementById("hamburger").classList.toggle("hide");
    document.getElementById("link_list").classList.toggle("hide");
}

// window.addEventListener('resize', function () {
// 	if (window.innerWidth < 425) {
//         document.getElementById("hamburger").classList.toggle("hide");
//         document.getElementById("link_list").classList.toggle("hide");
//     }
// })

hamburger = document.getElementById("hamburger_btn")
hamburger.onclick = function() {
    console.log(hamburger.innerHTML)
    if(hamburger.innerHTML == '<i class="bx bx-x"></i>') {
        hamburger.innerHTML = '<i class="bx bx-menu"></i>';
    } else {
        hamburger.innerHTML = '<i class="bx bx-x"></i>';
    }
    document.getElementById("link_list").classList.toggle("hide");
}