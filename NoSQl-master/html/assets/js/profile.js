// JavaScript for User Profile Page

// Function to open a specific tab content
function openTab(evt, tabName) {
    // Declare variables
    var i, tabcontent, tablinks;

    // Hide all tab content
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    // Deactivate all tab links
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }

    // Show the specific tab content and activate the tab link
    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className += " active";
}

// Get the tab elements and add event listeners
document.addEventListener("DOMContentLoaded", function() {
    var tablinks = document.getElementsByClassName("tablinks");
    for (var i = 0; i < tablinks.length; i++) {
        tablinks[i].addEventListener("click", function(event) {
            openTab(event, this.getAttribute("data-tab"));
        });
    }
});
