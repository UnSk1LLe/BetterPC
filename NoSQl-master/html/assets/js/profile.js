function openTab(evt, tabName) {
    var i, tabcontent, tablinks;

    //Hide all tab content
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }

    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className += " active";
}

document.addEventListener("DOMContentLoaded", function() {
    var tablinks = document.getElementsByClassName("tablinks");
    for (var i = 0; i < tablinks.length; i++) {
        tablinks[i].addEventListener("click", function(event) {
            openTab(event, this.getAttribute("data-tab"));
        });
    }

    var items = document.getElementsByClassName('item-details');
    for (var i = 0; i < items.length; i++) {
        var price = parseFloat(items[i].getAttribute('data-price'));
        var amount = parseFloat(items[i].getAttribute('data-amount'));
        var totalPrice = price * amount;
        items[i].querySelector('.total-price').textContent = totalPrice.toFixed(2);
    }
});

function confirmCancel(productId) {
    if (confirm("Are you sure you want cancel the order with ID: " + productId)) {
        var cancelForm = document.createElement('form');
        cancelForm.method = 'post';
        cancelForm.action = '/cancelOrder';
        var input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'orderID';
        input.value = productId;
        cancelForm.appendChild(input);
        document.body.appendChild(cancelForm);
        cancelForm.submit();
    }
}

function toggleDetails(button) {
    const orderDetails = button.parentElement.nextElementSibling;
    if (orderDetails.style.display === "none" || orderDetails.style.display === "") {
        orderDetails.style.display = "block";
        button.textContent = "Hide Details";
    } else {
        orderDetails.style.display = "none";
        button.textContent = "Show Details";
    }
}
