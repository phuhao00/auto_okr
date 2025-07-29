'use client'

import { useState } from 'react'
import { GitBranch, Calendar, FileText, Download, Loader2 } from 'lucide-react'
import axios from 'axios'

interface ReportData {
  content: string
  type: 'daily' | 'weekly'
  date: string
}

export default function Home() {
  const [repoPath, setRepoPath] = useState('')
  const [reportType, setReportType] = useState<'daily' | 'weekly'>('daily')
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0])
  const [loading, setLoading] = useState(false)
  const [report, setReport] = useState<ReportData | null>(null)
  const [error, setError] = useState('')

  const generateReport = async () => {
    if (!repoPath.trim()) {
      setError('请输入仓库路径')
      return
    }

    setLoading(true)
    setError('')
    setReport(null)

    try {
      const response = await axios.post('/api/generate-report', {
        repoPath: repoPath.trim(),
        type: reportType,
        date: selectedDate
      })

      setReport(response.data)
    } catch (err: any) {
      setError(err.response?.data?.error || '生成报告时发生错误')
    } finally {
      setLoading(false)
    }
  }

  const downloadReport = () => {
    if (!report) return

    const blob = new Blob([report.content], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${report.type}-report-${report.date}.md`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  return (
    <div className="px-4 py-8">
      <div className="max-w-4xl mx-auto">
        {/* 标题区域 */}
        <div className="text-center mb-8">
          <h2 className="text-3xl font-bold text-gray-900 mb-4">Git 提交报告生成器</h2>
          <p className="text-lg text-gray-600">输入仓库路径，生成精美的日报或周报</p>
        </div>

        {/* 输入表单 */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-8">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 仓库路径输入 */}
            <div className="md:col-span-2">
              <label htmlFor="repoPath" className="block text-sm font-medium text-gray-700 mb-2">
                <GitBranch className="inline w-4 h-4 mr-1" />
                仓库路径
              </label>
              <input
                type="text"
                id="repoPath"
                value={repoPath}
                onChange={(e) => setRepoPath(e.target.value)}
                placeholder="例如: /path/to/your/repo 或 C:\\path\\to\\repo"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            {/* 报告类型选择 */}
            <div>
              <label htmlFor="reportType" className="block text-sm font-medium text-gray-700 mb-2">
                <FileText className="inline w-4 h-4 mr-1" />
                报告类型
              </label>
              <select
                id="reportType"
                value={reportType}
                onChange={(e) => setReportType(e.target.value as 'daily' | 'weekly')}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="daily">日报</option>
                <option value="weekly">周报</option>
              </select>
            </div>

            {/* 日期选择 */}
            <div>
              <label htmlFor="date" className="block text-sm font-medium text-gray-700 mb-2">
                <Calendar className="inline w-4 h-4 mr-1" />
                目标日期
              </label>
              <input
                type="date"
                id="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          {/* 生成按钮 */}
          <div className="mt-6">
            <button
              onClick={generateReport}
              disabled={loading}
              className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-medium py-2 px-4 rounded-md transition duration-200 flex items-center justify-center"
            >
              {loading ? (
                <>
                  <Loader2 className="animate-spin w-4 h-4 mr-2" />
                  生成中...
                </>
              ) : (
                '生成报告'
              )}
            </button>
          </div>
        </div>

        {/* 错误信息 */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
            <div className="text-red-800">{error}</div>
          </div>
        )}

        {/* 报告结果 */}
        {report && (
          <div className="bg-white rounded-lg shadow-md p-6">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xl font-semibold text-gray-900">
                {report.type === 'daily' ? '日报' : '周报'} - {report.date}
              </h3>
              <button
                onClick={downloadReport}
                className="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md transition duration-200 flex items-center"
              >
                <Download className="w-4 h-4 mr-2" />
                下载
              </button>
            </div>
            <div className="bg-gray-50 rounded-md p-4 overflow-auto">
              <pre className="whitespace-pre-wrap text-sm text-gray-800">{report.content}</pre>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}