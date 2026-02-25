import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Router } from '../types'

interface RouterState {
  routers: Router[]
  selectedRouter: Router | null
  isLoading: boolean
  
  // Actions
  setRouters: (routers: Router[]) => void
  setSelectedRouter: (router: Router | null) => void
  addRouter: (router: Router) => void
  updateRouter: (id: string | number, router: Partial<Router>) => void
  removeRouter: (id: string | number) => void
  setLoading: (loading: boolean) => void
}

export const useRouterStore = create<RouterState>()(
  persist(
    (set, get) => ({
      routers: [],
      selectedRouter: null,
      isLoading: false,

      setRouters: (routers) => set({ routers }),
      getRouter: (id: string | number) => {
        const state = get()
        return state.routers.find((r) => r.id == id)
      },

      setSelectedRouter: (router) => set({ selectedRouter: router }),

      addRouter: (router) =>
        set((state) => ({
          routers: [...state.routers, router],
        })),

      updateRouter: (id: string | number, updatedRouter) =>
        set((state) => ({
          routers: state.routers.map((r) =>
            r.id == id ? { ...r, ...updatedRouter } : r
          ),
          selectedRouter:
            state.selectedRouter?.id == id
              ? { ...state.selectedRouter, ...updatedRouter }
              : state.selectedRouter,
        })),

      removeRouter: (id: string | number) =>
        set((state) => ({
          routers: state.routers.filter((r) => r.id != id),
          selectedRouter:
            state.selectedRouter?.id == id ? null : state.selectedRouter,
        })),

      setLoading: (isLoading) => set({ isLoading }),
    }),
    {
      name: 'mikhmon-router-v2',
      partialize: (state) => ({ selectedRouter: state.selectedRouter }),
    }
  )
)
