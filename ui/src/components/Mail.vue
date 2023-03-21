<script setup>
import { ref } from 'vue'
import { hasAttachment } from '../backend'
import {
    TrashIcon,
    EnvelopeIcon,
    ChevronDownIcon,
    ChevronUpIcon,
    EllipsisVerticalIcon,
    PaperClipIcon,
    XMarkIcon,
    CheckCircleIcon
} from '@heroicons/vue/24/outline'
import {
    Dialog,
    DialogPanel,
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue'

const props = defineProps(['message'])
const emits = defineEmits(['ondelete', 'onunread', 'onnext', 'onprev'])
const viewFrame = ref(null)
const showTip = ref(false)
const tipText = ref('')

async function onMarkUnread() {
    emits('onunread', props.message.id)
}

async function onDeleteMail() {
    emits('ondelete', props.message.id)
}

async function onNext() {
    emits('onnext', props.message.id)
}

async function onPrev() {
    emits('onprev', props.message.id)
}

async function onCopyContent() {
    navigator.clipboard.writeText(viewFrame.value.contentWindow.document.body.innerHTML);
    showTipText('Copy done')
}

function showTipText(text) {
    tipText.value = text
    showTip.value = true
    setTimeout(() => { showTip.value = false }, 5000)
}

async function onViewOriginal() {
    viewFrame.value.src = `/api/raw/${props.message.id}.eml`
}

function getAttachments() {
    return JSON.parse(props.message.attachments) || []
}

function onDownloadAttachment(path, name) {
    let link = document.createElement('A')
    link.href = `/api/raw/${path}`
    link.setAttribute('target', "_blank")
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
}

</script>
<template>
    <section aria-labelledby="message-heading" class="flex h-full min-w-0 flex-1 flex-col overflow-hidden xl:order-last">
        <div aria-live="assertive" class="pointer-events-none fixed inset-0 flex items-end px-4 py-6 sm:items-start sm:p-6">
            <div class="flex w-full flex-col items-center space-y-4 sm:items-end">
                <!-- Notification panel, dynamically insert this into the live region when it needs to be displayed -->
                <transition enter-active-class="transform ease-out duration-300 transition"
                    enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
                    enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
                    leave-active-class="transition ease-in duration-100" leave-from-class="opacity-100"
                    leave-to-class="opacity-0">
                    <div v-if="showTip"
                        class="pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg bg-white shadow-lg ring-1 ring-black ring-opacity-5">
                        <div class="p-4">
                            <div class="flex items-start">
                                <div class="flex-shrink-0">
                                    <CheckCircleIcon class="h-6 w-6 text-green-400" aria-hidden="true" />
                                </div>
                                <div class="ml-3 w-0 flex-1 pt-0.5">
                                    <p class="text-sm font-medium text-gray-900">{{ tipText }}</p>
                                </div>
                                <div class="ml-4 flex flex-shrink-0">
                                    <button type="button" @click="showTip = false"
                                        class="inline-flex rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                                        <span class="sr-only">Close</span>
                                        <XMarkIcon class="h-5 w-5" aria-hidden="true" />
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </transition>
            </div>
        </div>
        <!-- Top section -->
        <div class="flex-shrink-0 border-b border-gray-200 bg-white">
            <!-- Toolbar-->
            <div class="flex h-16 flex-col justify-center">
                <div class="px-4 sm:px-6 lg:px-8">
                    <div class="flex justify-between py-3">
                        <!-- Left buttons -->
                        <div>
                            <div class="isolate inline-flex rounded-md shadow-sm sm:space-x-3 sm:shadow-none">
                                <span class="inline-flex sm:shadow-sm">
                                    <button type="button"
                                        class="relative inline-flex items-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10"
                                        @click="onMarkUnread">
                                        <EnvelopeIcon class="-ml-0.5 h-5 w-5 text-gray-400" aria-hidden="true" />
                                        Mark unread
                                    </button>
                                </span>
                                <span class="inline-flex sm:shadow-sm">
                                    <button type="button" @click="onDeleteMail"
                                        class="relative inline-flex items-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10">
                                        <TrashIcon class="-ml-0.5 h-5 w-5 text-gray-400" aria-hidden="true" />
                                        Delete
                                    </button>
                                </span>
                            </div>
                        </div>

                        <!-- Right buttons -->
                        <nav aria-label="Pagination">
                            <span class="isolate inline-flex rounded-md shadow-sm">
                                <a href="#" @click="onPrev"
                                    class="relative -ml-px inline-flex items-center rounded-l-md bg-white px-3 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10">
                                    <span class="sr-only">Previous</span>
                                    <ChevronUpIcon class="h-5 w-5" aria-hidden="true" />
                                </a>
                                <a href="#" @click="onNext"
                                    class="relative inline-flex items-center rounded-r-md bg-white px-3 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10">
                                    <span class="sr-only">Next</span>
                                    <ChevronDownIcon class="h-5 w-5" aria-hidden="true" />
                                </a>
                            </span>
                        </nav>
                    </div>
                </div>
            </div>
            <!-- Message header -->
        </div>

        <div class="min-h-0 flex-1 overflow-y-auto">
            <div class="bg-white pt-5 pb-6 shadow">
                <div class="px-4 sm:flex sm:items-baseline sm:justify-between sm:px-6 lg:px-8">
                    <div class="sm:w-0 sm:flex-1">
                        <h1 id="message-heading" class="text-lg font-medium text-gray-900">{{ message.subject }}</h1>
                        <p>
                            <span class="mt-1 truncate text-sm text-gray-500 mx-4">From</span>
                            <span class="mt-1 truncate text-sm text-gray-700">{{ message.from }}</span>
                        </p>
                        <p>
                            <span class="mt-1 truncate text-sm text-gray-500 mx-4">To</span>
                            <span class="mt-1 truncate text-sm text-gray-700">{{ message.to }}</span>
                        </p>
                    </div>

                    <div class="mt-4 flex items-center justify-between sm:mt-0 sm:ml-6 sm:flex-shrink-0 sm:justify-start">
                        <span class="inline-flex items-center rounded-full px-3 py-0.5 text-sm font-medium "
                            :class="message.status_style">{{
                                message.status }}</span>
                        <Menu v-if="hasAttachment(message)" as="div" class="relative ml-3 inline-block text-left">
                            <div>
                                <MenuButton
                                    class="-my-2 flex items-center rounded-full bg-white p-2 text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-600">
                                    <span class="sr-only">Download attachments</span>
                                    <PaperClipIcon class="h-5 w-5" aria-hidden="true" />
                                </MenuButton>
                            </div>

                            <transition enter-active-class="transition ease-out duration-100"
                                enter-from-class="transform opacity-0 scale-95"
                                enter-to-class="transform opacity-100 scale-100"
                                leave-active-class="transition ease-in duration-75"
                                leave-from-class="transform opacity-100 scale-100"
                                leave-to-class="transform opacity-0 scale-95">
                                <MenuItems
                                    class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                    <div class="py-1">
                                        <MenuItem v-for="att in getAttachments(message)" :key="att.id">
                                        <a href="#" class="text-gray-700 flex justify-between px-4 py-2 text-sm"
                                            @click="onDownloadAttachment(att.path, att.name)">
                                            <span>{{ att.name }}</span> <span class="text-gray-500">{{ (att.size / (1024 *
                                                1024)).toFixed(2)
                                            }}MB</span>
                                        </a>
                                        </MenuItem>
                                    </div>
                                </MenuItems>
                            </transition>
                        </Menu>
                        <Menu as="div" class="relative ml-3 inline-block text-left">
                            <div>
                                <MenuButton
                                    class="-my-2 flex items-center rounded-full bg-white p-2 text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-600">
                                    <span class="sr-only">Open options</span>
                                    <EllipsisVerticalIcon class="h-5 w-5" aria-hidden="true" />
                                </MenuButton>
                            </div>

                            <transition enter-active-class="transition ease-out duration-100"
                                enter-from-class="transform opacity-0 scale-95"
                                enter-to-class="transform opacity-100 scale-100"
                                leave-active-class="transition ease-in duration-75"
                                leave-from-class="transform opacity-100 scale-100"
                                leave-to-class="transform opacity-0 scale-95">
                                <MenuItems
                                    class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                                    <div class="py-1">
                                        <MenuItem v-slot="{ active }">
                                        <button type="button" @click="onCopyContent"
                                            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'flex w-full justify-between px-4 py-2 text-sm']">
                                            <span>Copy content</span>
                                        </button>
                                        </MenuItem>
                                        <MenuItem v-slot="{ active }">
                                        <a href="#"
                                            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'flex justify-between px-4 py-2 text-sm']"
                                            @click="onViewOriginal">
                                            <span>View original</span>
                                        </a>
                                        </MenuItem>
                                    </div>
                                </MenuItems>
                            </transition>
                        </Menu>
                    </div>
                </div>
            </div>
            <!-- Mail content -->
            <div class="space-y-2 py-4 sm:space-y-4 sm:px-6 lg:px-8 min-w-full">
                <div class="bg-white px-4 py-6 shadow sm:rounded-lg sm:px-6 min-h-screen">
                    <iframe ref="viewFrame" :src="`/api/render/${message.id}`" class="min-w-full min-h-screen"></iframe>
                </div>
            </div>
        </div>
    </section>
</template>