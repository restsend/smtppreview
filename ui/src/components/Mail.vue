<script setup>
import {
    TrashIcon,
    EnvelopeIcon,
    ChevronDownIcon,
    ChevronUpIcon,
    EllipsisVerticalIcon,
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
const emits = defineEmits(['ondelete', 'onunread'])

async function onMarkUnread() {
    emits('onunread', props.message.id)
}

async function onDeleteMail() {
    emits('ondelete', props.message.id)
}

</script>
<template>
    <section aria-labelledby="message-heading" class="flex h-full min-w-0 flex-1 flex-col overflow-hidden xl:order-last">
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
                                <a href="#"
                                    class="relative inline-flex items-center rounded-l-md bg-white px-3 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10">
                                    <span class="sr-only">Next</span>
                                    <ChevronUpIcon class="h-5 w-5" aria-hidden="true" />
                                </a>
                                <a href="#"
                                    class="relative -ml-px inline-flex items-center rounded-r-md bg-white px-3 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:z-10 hover:bg-gray-50 focus:z-10">
                                    <span class="sr-only">Previous</span>
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
                                        <button type="button"
                                            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'flex w-full justify-between px-4 py-2 text-sm']">
                                            <span>Copy content</span>
                                        </button>
                                        </MenuItem>
                                        <MenuItem v-slot="{ active }">
                                        <a href="#"
                                            :class="[active ? 'bg-gray-100 text-gray-900' : 'text-gray-700', 'flex justify-between px-4 py-2 text-sm']">
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
                    <iframe :src="`/api/render/${message.id}`" class="min-w-full min-h-screen"></iframe>
                </div>
            </div>
        </div>
    </section>
</template>