
import './styles/colors.scss'
import './styles/typography.scss'
import './styles/document.scss'
import './styles/layout.scss'
import './styles/form.scss'
import './styles/content.scss'
import './styles/results.scss'

import Alpine from 'alpinejs'

interface FormattedResult {
  lbl: string;
  val: string | FormattedResult[];
};


function whoisLookup() {
  return {
    domain: '',
    result: '',
    error: '',
    formattedResult: [] as FormattedResult[],
    showFormatted: true,
    toggleFormat() {
      this.showFormatted = !this.showFormatted;
    },
    lookup() {
      fetch(`/${this.domain}`)
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.json();
        })
        .then(data => {
          this.result = JSON.stringify(data, null, 2);
          this.formattedResult = this.formatResult(data);
          this.error = '';
          console.log(this.formatResult(data));
        })
        .catch(error => {
          this.error = 'Failed to fetch WHOIS data';
          this.result = '';
          this.formattedResult = [];
          console.error('There has been a problem with your fetch operation:', error);
        });
    },
    formatResult(data: any) {
      const result: FormattedResult[] = [];
      Object.keys(data).forEach(key => {
        const sectionResults: FormattedResult[] = [];
        const label = key.replace(/_/g, ' ');
        const values = data[key];
        Object.keys(values).forEach(lbl => {
          const val = values[lbl];
          sectionResults.push({lbl, val});
        });
        result.push({ lbl: label, val: sectionResults });
      });
      return result;
    }
  };
}


Alpine.data('whoisLookup', whoisLookup);

window.Alpine = Alpine

Alpine.start()
