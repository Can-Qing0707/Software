<template>
  <div class="app-container">
    <template v-if="isStudent">
      <el-tabs v-model="activeTab" type="border-card" @tab-change="onTabChange">
        <el-tab-pane label="未完成任务" name="unfinished">
          <div class="filter-container">
            <el-select v-model="unfinQuery.course_id" placeholder="选择课程" clearable style="width:180px" @change="resetUnfinPage">
              <el-option v-for="c in courses" :key="c.id" :label="c.name" :value="c.id" />
            </el-select>
            <el-input v-model="unfinQuery.keyword" placeholder="搜索任务" clearable style="width:200px" @keyup.enter="resetUnfinPage" @clear="resetUnfinPage" />
            <el-button type="primary" @click="resetUnfinPage">搜索</el-button>
            <el-button type="success" @click="$router.push('/submission/upload')">提交成果</el-button>
          </div>
          <el-table :data="pagedUnfinishedTasks" stripe size="small">
            <el-table-column label="" width="70">
              <template #default>
                <el-tag type="warning" size="small">待提交</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="title" label="任务名称" min-width="160" />
            <el-table-column prop="course_name" label="所属课程" width="140" />
            <el-table-column prop="deadline" label="截止时间" width="160" />
            <el-table-column label="操作" width="160">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="$router.push('/task/' + row.id)">查看详情</el-button>
                <el-button type="primary" size="small" @click="$router.push('/submission/upload/' + row.id)">提交</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!filteredUnfinished.length" description="暂无未完成任务" />
          <div class="pagination-wrap" v-if="filteredUnfinished.length > unfinPageSize">
            <el-pagination
              v-model:current-page="unfinPage"
              :page-size="unfinPageSize"
              :total="filteredUnfinished.length"
              layout="prev, pager, next"
              background
              small
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="已提交记录" name="submitted">
          <div class="filter-container">
            <el-select v-model="subQuery.task_id" placeholder="选择任务" clearable style="width:200px" @change="resetSubPage">
              <el-option v-for="t in tasks" :key="t.id" :label="t.title" :value="t.id" />
            </el-select>
            <el-select v-model="subQuery.status" placeholder="状态" clearable style="width:140px" @change="resetSubPage">
              <el-option label="已上传" value="uploaded" />
              <el-option label="解析中" value="parsing" />
              <el-option label="已解析" value="parsed" />
              <el-option label="已核查" value="verified" />
              <el-option label="已评分" value="evaluated" />
            </el-select>
            <el-button type="primary" @click="resetSubPage">搜索</el-button>
            <el-button type="success" @click="$router.push('/submission/upload')">提交成果</el-button>
          </div>
          <el-table :data="pagedSubmissions" stripe size="small">
            <el-table-column prop="task_title" label="任务" min-width="180" />
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="statusMap[row.status]?.type || 'info'" size="small">
                  {{ statusMap[row.status]?.label || row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="submit_time" label="提交时间" width="160" />
            <el-table-column label="评分结果" width="100">
              <template #default="{ row }">
                <template v-if="row.total_score != null">
                  <span style="font-weight:bold;color:#409eff">{{ row.total_score.toFixed(1) }} 分</span>
                </template>
                <template v-else>
                  <el-tag type="info" size="small">未评分</el-tag>
                </template>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="$router.push('/submission/' + row.id)">查看</el-button>
                <el-button type="warning" link size="small" @click="$router.push('/submission/upload/' + row.task_id)">重新提交</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!filteredSubmissions.length" description="暂无提交记录" />
          <div class="pagination-wrap" v-if="filteredSubmissions.length > subPageSize">
            <el-pagination
              v-model:current-page="subPage"
              :page-size="subPageSize"
              :total="filteredSubmissions.length"
              layout="prev, pager, next"
              background
              small
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </template>

    <template v-else>
      <div class="filter-container">
        <el-select v-model="subQuery.task_id" placeholder="选择任务" clearable @change="resetSubPage" style="width:200px">
          <el-option v-for="t in tasks" :key="t.id" :label="t.title" :value="t.id" />
        </el-select>
        <el-select v-model="subQuery.status" placeholder="状态" clearable @change="resetSubPage" style="width:140px">
          <el-option label="已上传" value="uploaded" />
          <el-option label="解析中" value="parsing" />
          <el-option label="已解析" value="parsed" />
          <el-option label="已核查" value="verified" />
          <el-option label="已评分" value="evaluated" />
        </el-select>
        <el-button type="primary" @click="resetSubPage">搜索</el-button>
      </div>
      <el-table :data="pagedSubmissions" stripe>
        <el-table-column prop="student_name" label="学生" width="120" />
        <el-table-column prop="task_title" label="任务" min-width="160" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'evaluated' ? 'success' : (statusMap[row.status]?.type || 'info')" size="small">
              {{ row.status === 'evaluated' ? '已评分' : (statusMap[row.status]?.label || row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="LLM评分" width="90">
          <template #default="{ row }">
            {{ row.llm_total != null ? row.llm_total.toFixed(1) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="教师评分" width="90">
          <template #default="{ row }">
            <span :style="{ color: row.teacher_total != null ? '#67c23a' : '#909399', fontWeight: row.teacher_total != null ? 'bold' : 'normal' }">
              {{ row.teacher_total != null ? row.teacher_total.toFixed(1) : '-' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="submit_time" label="提交时间" width="160" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push('/submission/' + row.id)">查看</el-button>
            <el-button type="warning" link @click="$router.push('/evaluation/score/' + row.id)">评分</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!filteredSubmissions.length" description="暂无数据" />
      <div class="pagination-wrap" v-if="filteredSubmissions.length > subPageSize">
        <el-pagination
          v-model:current-page="subPage"
          :page-size="subPageSize"
          :total="filteredSubmissions.length"
          layout="prev, pager, next"
          background
          small
        />
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getSubmissionList } from '@/api/submission'
import { getTaskList } from '@/api/task'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const isStudent = computed(() => userStore.role === 'student')

const activeTab = ref('unfinished')

const allTasks = ref([])
const allSubmissions = ref([])
const courses = ref([])

const unfinQuery = ref({ course_id: '', keyword: '' })
const subQuery = ref({ task_id: '', status: '' })

const unfinPage = ref(1)
const subPage = ref(1)
const unfinPageSize = 6
const subPageSize = 6

const statusMap = {
  uploaded: { label: '已上传', type: 'info' },
  parsing: { label: '解析中', type: 'warning' },
  parsed: { label: '已解析', type: '' },
  verified: { label: '已核查', type: 'success' },
  evaluated: { label: '已评分', type: 'success' }
}

const submittedTaskIds = computed(() => new Set(allSubmissions.value.map(s => s.task_id)))

const unfinishedTasks = computed(() => allTasks.value.filter(t => !submittedTaskIds.value.has(t.id)))

const filteredUnfinished = computed(() => {
  let list = unfinishedTasks.value
  const q = unfinQuery.value
  if (q.course_id) {
    list = list.filter(t => t.course_id === Number(q.course_id))
  }
  if (q.keyword) {
    const kw = q.keyword.toLowerCase()
    list = list.filter(t => t.title.toLowerCase().includes(kw))
  }
  return list
})

const filteredSubmissions = computed(() => {
  let list = allSubmissions.value
  const q = subQuery.value
  if (q.task_id) {
    list = list.filter(s => s.task_id === Number(q.task_id))
  }
  if (q.status) {
    list = list.filter(s => s.status === q.status)
  }
  return list
})

const pagedUnfinishedTasks = computed(() => {
  const start = (unfinPage.value - 1) * unfinPageSize
  return filteredUnfinished.value.slice(start, start + unfinPageSize)
})

const pagedSubmissions = computed(() => {
  const start = (subPage.value - 1) * subPageSize
  return filteredSubmissions.value.slice(start, start + subPageSize)
})

function resetUnfinPage() {
  unfinPage.value = 1
}

function resetSubPage() {
  subPage.value = 1
}

function onTabChange() {
  // reset page when switching tabs
}

onMounted(async () => {
  try {
    const [taskRes, subRes] = await Promise.all([
      getTaskList(),
      getSubmissionList()
    ])
    const tasks = taskRes.data || []
    const submissions = subRes.data || []
    allTasks.value = tasks
    allSubmissions.value = submissions
    const courseMap = {}
    tasks.forEach(t => {
      courseMap[t.course_id] = { id: t.course_id, name: t.course_name }
    })
    courses.value = Object.values(courseMap)
  } catch {
    // ignore
  }
})
</script>

<style scoped>
.filter-container {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}
</style>
