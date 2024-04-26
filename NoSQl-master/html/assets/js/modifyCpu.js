document.addEventListener("DOMContentLoaded", function() {
    const ecores = document.getElementById('ecores');
    const ecoresBase = document.getElementById('ecoresBase');
    const ecoresBoost = document.getElementById('ecoresBoost');

    function setDisabledState() {
        if (ecores.value === '0') {
            ecoresBase.disabled = true;
            ecoresBoost.disabled = true;
        } else {
            ecoresBase.disabled = false;
            ecoresBoost.disabled = false;
        }
    }

    setDisabledState();

    ecores.addEventListener('input', setDisabledState);
});