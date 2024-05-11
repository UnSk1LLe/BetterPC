document.addEventListener('DOMContentLoaded', function() {
    var productData = document.getElementById('productData');
    var inCart = productData.dataset.inCart === 'true';
    var amount = parseInt(productData.dataset.amount);

    var addToCartButton = document.getElementById('addToCartButton');
    var outOfStockMessage = document.getElementById('outOfStockMessage');
    var addedText = document.getElementById('addedText');
wa
    if (amount > 0) {
        if (inCart) {
            addToCartButton.style.display = 'none';
            addedText.style.display = 'block';
        } else {
            addToCartButton.style.display = 'block';
            addedText.style.display = 'none';

            addToCartButton.addEventListener('click', function() {
                if (amount > 0) {
                    addToCart();
                } else {
                    outOfStockMessage.style.display = 'block';
                }
            });
        }
    } else {
        outOfStockMessage.style.display = 'block';
        addToCartButton.style.display = 'none';
        addedText.style.display = 'none';
    }

    function addToCart() {
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/addProductToCart', true);
        xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
        xhr.onreadystatechange = function() {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                if (xhr.status === 200) {
                    alert('Product added to cart');
                    addToCartButton.style.display = 'none';
                    addedText.style.display = 'block';
                } else {
                    alert('Error adding product to cart');
                }
            }
        };
        xhr.send('addToCart=cpu' + productData.dataset.id);
    }
});
