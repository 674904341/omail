import React, { useState, useEffect } from "react"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"

function DisclosureModal() {
  const [open, setOpen] = useState(false)

  useEffect(() => {
    // 检查今天是否已经显示过弹窗
    const today = new Date().toDateString()
    const lastShown = localStorage.getItem("disclosureModalLastShown")

    if (lastShown !== today) {
      setOpen(true)
    }
  }, [])

  const handleDismiss = () => {
    setOpen(false)
  }

  const handleDismissToday = () => {
    const today = new Date().toDateString()
    localStorage.setItem("disclosureModalLastShown", today)
    setOpen(false)
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogContent>
        <AlertDialogTitle>重要提示</AlertDialogTitle>
        <AlertDialogDescription className="space-y-2">
          <div>❗接收到的邮件内容仅能保留10天</div>
          <div>❗随机生成的邮箱地址任何人都可以使用，请勿用于注册重要账号</div>
          <div>❗请勿发送包含敏感信息的邮件至该邮箱</div>
          <div>❗本服务不保证100%稳定可用，可能会出现无法接收邮件的情况</div>
          <div>❗请勿将该邮箱用于垃圾邮件发送等违法行为，否则后果自负</div>
          <div>❗使用本服务即表示您已知悉并同意以上事项</div>
        </AlertDialogDescription>
        <div className="flex gap-2 justify-end">
          <AlertDialogCancel onClick={handleDismiss}>
            知悉并同意
          </AlertDialogCancel>
          <AlertDialogAction onClick={handleDismissToday}>
            今日不再提示
          </AlertDialogAction>
        </div>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default DisclosureModal
