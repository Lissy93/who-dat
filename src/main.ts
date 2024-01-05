
import './styles/colors.scss'
import './styles/typography.scss'
import './styles/document.scss'
import './styles/layout.scss'
import './styles/form.scss'
import './styles/content.scss'

import Alpine from 'alpinejs'

function whoisLookup() {
  return {
      domain: '',
      result: '',
      error: '',
      lookup() {
        // fetch(`https://potential-bassoon-77r4qqr57g2x5xr-3000.app.github.dev/${this.domain}`)
        fetch(`/${this.domain}`)
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.json();
        })
        .then(data => {
          this.result = JSON.stringify(data, null, 2);
          this.error = '';
        })
        .catch(error => {
          this.error = 'Failed to fetch WHOIS data';
          this.result = '';
          console.error('There has been a problem with your fetch operation:', error);
        });
      }
  };
}

Alpine.data('whoisLookup', whoisLookup);

window.Alpine = Alpine

Alpine.start()
