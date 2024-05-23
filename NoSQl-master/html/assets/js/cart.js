function changeQuantity(index, action) {
    var quantityElement = document.getElementById('quantity' + index);
    var currentQuantity = parseInt(quantityElement.textContent);
    var cardElement = document.querySelector('.cpu-card[data-index="' + index + '"]');
    var maxAmount = parseInt(cardElement.getAttribute('data-max-amount'));
    var pricePerUnit = parseInt(cardElement.getAttribute('data-price'));
    var maxWarning = document.getElementById('maxWarning' + index);
    var priceElement = document.getElementById('price' + index);

    if (action === 'increase') {
        if (currentQuantity < maxAmount) {
            currentQuantity++;
        } else {
            maxWarning.style.display = 'block';
            return;
        }
    } else if (action === 'decrease' && currentQuantity > 1) {
        currentQuantity--;
        maxWarning.style.display = 'none';
    }

    quantityElement.textContent = currentQuantity;
    priceElement.textContent = currentQuantity * pricePerUnit;

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/updateCart', true);
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.send('index=' + index + '&quantity=' + currentQuantity);

    checkOrderButton();
}

function confirmDelete(productId) {
    if (confirm("Are you sure you want to delete the item from cart?")) {
        var deleteForm = document.createElement('form');
        deleteForm.method = 'post';
        deleteForm.action = '/deleteProductFromCart';
        var input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'deleteProduct';
        input.value = productId;
        deleteForm.appendChild(input);
        document.body.appendChild(deleteForm);
        deleteForm.submit();
    }
}

function checkOrderButton() {
    var createOrderButton = document.getElementById('createOrderButton');
    var allQuantitiesValid = true;
    var quantityElements = document.querySelectorAll('.cpu-card');

    quantityElements.forEach(function(cardElement) {
        var index = cardElement.getAttribute('data-index');
        var quantity = parseInt(document.getElementById('quantity' + index).textContent);
        var maxAmount = parseInt(cardElement.getAttribute('data-max-amount'));
        if (quantity > maxAmount) {
            allQuantitiesValid = false;
        }
    });

    if (allQuantitiesValid) {
        createOrderButton.disabled = false;
    } else {
        createOrderButton.disabled = true;
    }
}

// Initial check when page loads
window.onload = function() {
    checkOrderButton();
}