
<script setup>
import { ref } from 'vue'
import MailVue from './components/Mail.vue'
import { getSummary, listMails, markMailOpened, deleteMail, hasAttachment } from './backend'
import TimeAgo from 'javascript-time-ago'
// English.
import en from 'javascript-time-ago/locale/en'
import {
  Cog6ToothIcon,
  ArrowPathIcon,
  MagnifyingGlassIcon,
  PaperClipIcon
} from '@heroicons/vue/24/outline'

TimeAgo.addDefaultLocale(en)
const timeAgo = new TimeAgo('en-US')

const totalCount = ref(0)
const totalSize = ref(0)
const messages = ref([])
const message = ref(null)
const pos = ref(0)
const keyword = ref('')

onRefresh().then()

async function onRefresh() {
  let summary = await getSummary();
  totalCount.value = summary.totalCount
  totalSize.value = (summary.totalSize / (1024 * 1024)).toFixed(2)
  let data = await listMails()
  messages.value = data.items || []
}

async function onSelectMail(msg) {
  msg.opened = true
  await markMailOpened(msg.id, true)

  if (message.value) { message.value.selected = false }
  message.value = msg
  message.value.selected = true
}

function onKeywordChanged() {
  listMails(pos.value, keyword.value).then((data) => {
    messages.value = data.items || []
  })
}

async function onMarkUnread() {
  await markMailOpened(message.value.id, false)
  message.value.opened = false
}

async function onDeleteMail() {
  let idx = 0
  messages.value = messages.value.filter((e, i) => {
    if (e.id == message.value.id) { idx = i }
    return e.id != message.value.id
  })

  await deleteMail(message.value.id)
  if (idx >= messages.value.length) {
    idx = messages.value.length - 1
  }
  if (messages.value.length > 0) {
    await onSelectMail(messages.value[idx])
  } else {
    message.value = null
  }
}

async function onNext(msgid) {
  for (let idx = 0; idx < messages.value.length; idx++) {
    const e = messages.value[idx];
    if (e.id == msgid) {
      idx = idx + 1
      if (idx >= messages.value.length) {
        idx = messages.value.length - 1
      }
      await onSelectMail(messages.value[idx])
      return
    }
  }
}

async function onPrev(msgid) {
  for (let idx = 0; idx < messages.value.length; idx++) {
    const e = messages.value[idx];
    if (e.id == msgid) {
      idx = idx - 1
      if (idx < 0) { idx = 0 }
      await onSelectMail(messages.value[idx])
      return
    }
  }
}

</script>
<template>
  <div class="flex h-full flex-col">
    <!-- Top nav-->
    <header class="relative flex h-16 flex-shrink-0 items-center bg-white">
      <!-- Logo area -->
      <div class="absolute inset-y-0 left-0 lg:static lg:flex-shrink-0">
        <a href="#"
          class="flex h-16 w-16 items-center justify-center bg-cyan-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-600 lg:w-20">
          <img class="h-8 w-auto" src="/src/assets/logo.png" alt="SMTP Preview by restsend.com" />
        </a>
      </div>

      <!-- Desktop nav area -->
      <div class="hidden lg:flex lg:min-w-0 lg:flex-1 lg:items-center lg:justify-between">
        <div class="min-w-0 flex-1">
          <div class="relative max-w-2xl text-gray-400 focus-within:text-gray-500">
            <label for="desktop-search" class="sr-only">Search all inboxes</label>
            <input id="desktop-search" type="search" placeholder="Search all inboxes"
              class="block w-full border-transparent pl-12 text-gray-900 focus:border-transparent focus:ring-0 sm:text-sm"
              v-model="keyword" @keydown.enter="onKeywordChanged" />
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center justify-center pl-4">
              <MagnifyingGlassIcon class="h-5 w-5" aria-hidden="true" />
            </div>
          </div>
        </div>
        <div class="text-gray-700">
          <span class="mx-2">Total<span class="mx-1 font-bold">{{ totalCount }}</span>Messages</span>
          <span><span class="mx-1 font-bold">{{ totalSize }}</span>MB</span>
        </div>
        <div class="ml-10 flex flex-shrink-0 items-center space-x-10 pr-4">
          <div class="flex items-center space-x-8">
            <span class="inline-flex">
              <a href="#" @click="onRefresh" class="mx-1 rounded-full bg-white p-1 text-gray-400 hover:text-gray-500">
                <span class="sr-only">Refresh</span>
                <ArrowPathIcon class="h-6 w-6" aria-hidden="true" />
              </a>
              <a href="#" class="mx-1 rounded-full bg-white p-1 text-gray-400 hover:text-gray-500">
                <span class="sr-only">Settings</span>
                <Cog6ToothIcon class="h-6 w-6" aria-hidden="true" />
              </a>
            </span>
          </div>
        </div>
      </div>

    </header>

    <!-- Bottom section -->
    <div class="flex min-h-0 flex-1 overflow-hidden">
      <!-- Main area -->
      <main class="min-w-0 flex-1 border-t border-gray-200 xl:flex">
        <MailVue v-if="message" :message="message" @onunread="onMarkUnread" @ondelete="onDeleteMail" @onnext="onNext"
          @onprev="onPrev" />

        <!-- Message list-->
        <aside class="hidden xl:order-first xl:block xl:flex-shrink-0">
          <div class="relative flex h-full w-96 flex-col border-r border-gray-200 bg-gray-100">
            <div class="flex-shrink-0">
              <div class="flex h-16 flex-col justify-center bg-white px-6">
                <div class="flex items-baseline space-x-3">
                  <h2 class="text-lg font-medium text-gray-900">Inbox</h2>
                  <p class="text-sm font-medium text-gray-500">{{ messages.length }} messages</p>
                </div>
              </div>
              <div class="border-t border-b border-gray-200 bg-gray-50 px-6 py-2 text-sm font-medium text-gray-500">
                Sorted by date</div>
            </div>
            <nav aria-label="Message list" class="min-h-0 flex-1 overflow-y-auto">
              <ul role="list" class="divide-y divide-gray-200 border-b border-gray-200">
                <li v-for="msg in messages" :key="msg.id"
                  class="relative py-5 px-6 focus-within:ring-2 focus-within:ring-inset focus-within:ring-blue-600 hover:bg-cyan-50"
                  :class="[msg.selected ? 'bg-cyan-50' : 'bg-white']">
                  <a href="#" @click="onSelectMail(msg)">
                    <div class="flex justify-between space-x-3">
                      <div class="min-w-0 flex-1" :class="msg.opened ? '' : 'font-bold'">
                        <a :href="msg.href" class="block focus:outline-none">
                          <span class="absolute inset-0" aria-hidden="true" />
                          <p class="truncate text-sm font-medium text-gray-900">{{ msg.from }}</p>
                          <p class="truncate text-sm text-gray-500">{{ msg.subject }}</p>
                        </a>
                      </div>
                      <div class="flex justify-between">
                        <time :datetime="msg.createdAt" class="flex-shrink-0 whitespace-nowrap text-sm text-gray-500">{{
                          timeAgo.format(new Date(msg.createdAt)) }}</time>
                        <PaperClipIcon v-show="hasAttachment(msg)" class="mx-1 w-4 h-4"></PaperClipIcon>
                      </div>
                    </div>
                    <div v-if="msg.textBody" class="mt-1">
                      <p class="text-sm text-gray-600 line-clamp-2">{{ msg.textBody.substr(0, 128) }}</p>
                    </div>
                  </a>
                </li>
              </ul>
            </nav>
          </div>
        </aside>
      </main>
    </div>
  </div>
</template>
