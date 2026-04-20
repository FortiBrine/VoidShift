export const useNotification = () => {
  const visible = useState<boolean>('notification-visible', () => false)
  const message = useState<string>('notification-message', () => '')
  const color = useState<'success' | 'error' | 'info' | 'warning'>('notification-color', () => 'info')
  const timeout = 3500

  const show = (text: string, nextColor: 'success' | 'error' | 'info' | 'warning') => {
    message.value = text
    color.value = nextColor
    visible.value = true
  }

  const showSuccess = (text: string) => show(text, 'success')
  const showError = (text: string) => show(text, 'error')
  const showInfo = (text: string) => show(text, 'info')

  return {
    visible,
    message,
    color,
    timeout,
    show,
    showSuccess,
    showError,
    showInfo,
  }
}
