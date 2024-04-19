document.addEventListener("DOMContentLoaded", function() {
    document.getElementById('manufacturer').value = '{{ .Manufacturer }}';
    document.getElementById('freeMult').value = '{{ .ClockFrequency.FreeMultiplier }}';
    document.getElementById('type').value = '{{ .Ram.Type }}';
    document.getElementById('pcie').value = '{{ .PciE }}';

    const ecores = document.getElementById('ecores');
    const ecoresBase = document.getElementById('ecoresBase');
    const ecoresBoost = document.getElementById('ecoresBoost');

    if (ecores.value === 0) {
        ecoresBase.disabled = true;
        ecoresBoost.disabled = true;
    } else {
        ecoresBase.disabled = false;
        ecoresBoost.disabled = false;
    }
    ecores.addEventListener('input', function() {
        if (ecores.value > 0) {
            ecoresBase.disabled = false;
            ecoresBoost.disabled = false;
        } else {
            ecoresBase.disabled = true;
            ecoresBoost.disabled = true;
        }
    });
});