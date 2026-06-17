import { ref } from 'vue'

const msg = ref('')
const type = ref('success')

let timer = null

export function useToast() {
  function show(m, t = 'success') {
    msg.value = m; type.value = t
    clearTimeout(timer)
    timer = setTimeout(() => msg.value = '', 3000)
  }
  return { msg, type, success: (m) => show(m, 'success'), error: (m) => show(m, 'error'), warning: (m) => show(m, 'warning') }
}
