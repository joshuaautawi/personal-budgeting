import { useMemo, useState } from 'react'
import type { CategoryType, Id } from '../lib/types'
import { useAppStore } from '../store/AppStore'

export function Categories() {
  const { state, addCategory, updateCategory, deleteCategory } = useAppStore()

  const [type, setType] = useState<CategoryType>('expense')
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')

  const [editingId, setEditingId] = useState<Id | null>(null)
  const [editName, setEditName] = useState('')
  const [editDescription, setEditDescription] = useState('')

  const income = useMemo(() => state.categories.filter((c) => c.type === 'income'), [state.categories])
  const expense = useMemo(() => state.categories.filter((c) => c.type === 'expense'), [state.categories])

  return (
    <div className="grid">
      <section className="card">
        <div className="card__header">
          <div className="card__title">Create Category</div>
          <div className="card__hint">Income and expense categories are reusable across months.</div>
        </div>

        <div className="form">
          <label className="field">
            <div className="field__label">Type</div>
            <select className="select" value={type} onChange={(e) => setType(e.target.value as CategoryType)}>
              <option value="expense">Expense</option>
              <option value="income">Income</option>
            </select>
          </label>

          <label className="field">
            <div className="field__label">Name</div>
            <input className="input" value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. Groceries" />
          </label>

          <label className="field field--full">
            <div className="field__label">Description (optional)</div>
            <input
              className="input"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Any notes about this category"
            />
          </label>

          <div className="form__actions">
            <button
              className="btn"
              onClick={async () => {
                const res = await addCategory({ type, name, description })
                if (!res.ok) return
                setName('')
                setDescription('')
              }}
            >
              Add category
            </button>
          </div>
        </div>
      </section>

      <section className="card">
        <div className="card__header">
          <div className="card__title">Expense Categories</div>
          <div className="card__hint">Used for budgets and expense transactions.</div>
        </div>

        <CategoryList
          ids={expense.map((c) => c.id)}
          editingId={editingId}
          startEdit={(id) => {
            const cat = state.categories.find((c) => c.id === id)
            if (!cat) return
            setEditingId(id)
            setEditName(cat.name)
            setEditDescription(cat.description ?? '')
          }}
          cancelEdit={() => setEditingId(null)}
          saveEdit={async (id) => {
            const res = await updateCategory(id, { name: editName, description: editDescription })
            if (!res.ok) return
            setEditingId(null)
          }}
          remove={async (id) => {
            await deleteCategory(id)
          }}
          editName={editName}
          setEditName={setEditName}
          editDescription={editDescription}
          setEditDescription={setEditDescription}
          getCategory={(id) => state.categories.find((c) => c.id === id)}
        />
      </section>

      <section className="card">
        <div className="card__header">
          <div className="card__title">Income Categories</div>
          <div className="card__hint">Used for income transactions.</div>
        </div>

        <CategoryList
          ids={income.map((c) => c.id)}
          editingId={editingId}
          startEdit={(id) => {
            const cat = state.categories.find((c) => c.id === id)
            if (!cat) return
            setEditingId(id)
            setEditName(cat.name)
            setEditDescription(cat.description ?? '')
          }}
          cancelEdit={() => setEditingId(null)}
          saveEdit={async (id) => {
            const res = await updateCategory(id, { name: editName, description: editDescription })
            if (!res.ok) return
            setEditingId(null)
          }}
          remove={async (id) => {
            await deleteCategory(id)
          }}
          editName={editName}
          setEditName={setEditName}
          editDescription={editDescription}
          setEditDescription={setEditDescription}
          getCategory={(id) => state.categories.find((c) => c.id === id)}
        />
      </section>
    </div>
  )
}

function CategoryList(props: {
  ids: Id[]
  editingId: Id | null
  startEdit: (id: Id) => void
  cancelEdit: () => void
  saveEdit: (id: Id) => void
  remove: (id: Id) => void
  editName: string
  setEditName: (v: string) => void
  editDescription: string
  setEditDescription: (v: string) => void
  getCategory: (id: Id) => { id: Id; name: string; description?: string } | undefined
}) {
  if (props.ids.length === 0) return <div className="empty">No categories yet.</div>

  return (
    <div className="list">
      {props.ids.map((id) => {
        const cat = props.getCategory(id)
        if (!cat) return null
        const isEditing = props.editingId === id
        return (
          <div key={id} className="row">
            <div className="row__main">
              {isEditing ? (
                <div className="row__edit">
                  <input className="input" value={props.editName} onChange={(e) => props.setEditName(e.target.value)} />
                  <input
                    className="input"
                    value={props.editDescription}
                    onChange={(e) => props.setEditDescription(e.target.value)}
                    placeholder="Description (optional)"
                  />
                </div>
              ) : (
                <>
                  <div className="row__title">{cat.name}</div>
                  {cat.description ? <div className="row__sub">{cat.description}</div> : null}
                </>
              )}
            </div>

            <div className="row__actions">
              {isEditing ? (
                <>
                  <button className="btn btn--ghost" onClick={props.cancelEdit}>
                    Cancel
                  </button>
                  <button className="btn" onClick={() => props.saveEdit(id)}>
                    Save
                  </button>
                </>
              ) : (
                <>
                  <button className="btn btn--ghost" onClick={() => props.startEdit(id)}>
                    Edit
                  </button>
                  <button className="btn btn--danger" onClick={() => props.remove(id)}>
                    Delete
                  </button>
                </>
              )}
            </div>
          </div>
        )
      })}
    </div>
  )
}


