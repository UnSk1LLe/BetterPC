//Multiplies price by quantity. used onload
function updatePriceAndQuantity(index) {
    let quantityElement = document.getElementById('quantity' + index);
    let currentQuantity = parseInt(quantityElement.textContent);
    let cardElement = document.querySelector('.product-card[data-index="' + index + '"]');
    let pricePerUnit = parseInt(cardElement.getAttribute('data-price'));
    let priceElement = document.getElementById('price' + index);

    priceElement.textContent = currentQuantity * pricePerUnit;
}

//Changes quantity and price
function changeQuantity(index, action) {
    let quantityElement = document.getElementById('quantity' + index);
    let currentQuantity = parseInt(quantityElement.textContent);
    let cardElement = document.querySelector('.product-card[data-index="' + index + '"]');
    let maxAmount = parseInt(cardElement.getAttribute('data-max-amount'));
    let pricePerUnit = parseInt(cardElement.getAttribute('data-price'));
    let maxWarning = document.getElementById('maxWarning' + index);
    let priceElement = document.getElementById('price' + index);

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

    if (allQuantitiesValid && quantityElements) {
        createOrderButton.disabled = false;
    } else {
        createOrderButton.disabled = true;
    }
}

window.onload = function() {
    checkOrderButton();
    let productCards = document.querySelectorAll('.product-card');
    productCards.forEach(function(card) {
        let index = card.getAttribute('data-index');
        updatePriceAndQuantity(index);
    });
}

function showProduct(productType, productID) {
    window.location.href = `/showProduct?productType=${encodeURIComponent(productType)}&productID=${encodeURIComponent(productID)}`;
}