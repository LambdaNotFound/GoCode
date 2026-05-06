import { useState, useEffect, useCallback } from "react";

const API = "";

const CATEGORIES = ["Food", "Travel", "Software", "Hardware", "Marketing", "Other"];

function formatCurrency(amount) {
  return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(amount);
}

export default function App() {
  const [expenses, setExpenses] = useState([]);
  const [summary, setSummary] = useState([]);
  const [filterCategory, setFilterCategory] = useState("");
  const [error, setError] = useState("");
  const [form, setForm] = useState({
    amount: "",
    category: "Food",
    description: "",
    date: new Date().toISOString().split("T")[0],
  });

  const fetchExpenses = useCallback(async () => {
    const url = filterCategory
      ? `/expenses?category=${encodeURIComponent(filterCategory)}`
      : "/expenses";
    const res = await fetch(url);
    if (!res.ok) return;
    setExpenses(await res.json());
  }, [filterCategory]);

  const fetchSummary = async () => {
    const res = await fetch("/expenses/summary");
    if (!res.ok) return;
    setSummary(await res.json());
  };

  useEffect(() => {
    fetchExpenses();
    fetchSummary();
  }, [fetchExpenses]);

  const handleCreate = async (e) => {
    e.preventDefault();
    setError("");
    const res = await fetch("/expenses", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...form, amount: parseFloat(form.amount) }),
    });
    if (!res.ok) {
      const body = await res.json();
      setError(body.error || "Failed to create expense");
      return;
    }
    setForm({ amount: "", category: "Food", description: "", date: new Date().toISOString().split("T")[0] });
    fetchExpenses();
    fetchSummary();
  };

  const handleDelete = async (id) => {
    await fetch(`/expenses/${id}`, { method: "DELETE" });
    fetchExpenses();
    fetchSummary();
  };

  return (
    <div style={styles.page}>
      <h1 style={styles.title}>Expense Tracker</h1>

      {/* Summary */}
      {summary.length > 0 && (
        <div style={styles.card}>
          <h2 style={styles.sectionTitle}>Summary by Category</h2>
          <div style={styles.summaryGrid}>
            {summary.map((s) => (
              <div key={s.category} style={styles.summaryItem}>
                <div style={styles.summaryCategory}>{s.category}</div>
                <div style={styles.summaryTotal}>{formatCurrency(s.total)}</div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Add Expense Form */}
      <div style={styles.card}>
        <h2 style={styles.sectionTitle}>Add Expense</h2>
        <form onSubmit={handleCreate} style={styles.form}>
          <input
            style={styles.input}
            type="number"
            placeholder="Amount"
            min="0.01"
            step="0.01"
            value={form.amount}
            onChange={(e) => setForm({ ...form, amount: e.target.value })}
            required
          />
          <select
            style={styles.input}
            value={form.category}
            onChange={(e) => setForm({ ...form, category: e.target.value })}
          >
            {CATEGORIES.map((c) => <option key={c}>{c}</option>)}
          </select>
          <input
            style={styles.input}
            type="text"
            placeholder="Description (optional)"
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
          />
          <input
            style={styles.input}
            type="date"
            value={form.date}
            onChange={(e) => setForm({ ...form, date: e.target.value })}
            required
          />
          <button type="submit" style={styles.button}>Add</button>
        </form>
        {error && <p style={styles.error}>{error}</p>}
      </div>

      {/* Expense List */}
      <div style={styles.card}>
        <div style={styles.listHeader}>
          <h2 style={styles.sectionTitle}>Expenses</h2>
          <select
            style={{ ...styles.input, width: "auto" }}
            value={filterCategory}
            onChange={(e) => setFilterCategory(e.target.value)}
          >
            <option value="">All Categories</option>
            {CATEGORIES.map((c) => <option key={c}>{c}</option>)}
          </select>
        </div>

        {expenses.length === 0 ? (
          <p style={styles.empty}>No expenses yet.</p>
        ) : (
          <table style={styles.table}>
            <thead>
              <tr>
                {["Date", "Category", "Description", "Amount", ""].map((h) => (
                  <th key={h} style={styles.th}>{h}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {expenses.map((e) => (
                <tr key={e.id} style={styles.tr}>
                  <td style={styles.td}>{e.date}</td>
                  <td style={styles.td}>{e.category}</td>
                  <td style={styles.td}>{e.description || "—"}</td>
                  <td style={{ ...styles.td, fontWeight: 600 }}>{formatCurrency(e.amount)}</td>
                  <td style={styles.td}>
                    <button onClick={() => handleDelete(e.id)} style={styles.deleteBtn}>✕</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

const styles = {
  page: { maxWidth: 800, margin: "0 auto", padding: "2rem 1rem", fontFamily: "system-ui, sans-serif", color: "#111" },
  title: { fontSize: "1.75rem", fontWeight: 700, marginBottom: "1.5rem" },
  card: { background: "#fff", border: "1px solid #e5e7eb", borderRadius: 8, padding: "1.25rem", marginBottom: "1.25rem" },
  sectionTitle: { fontSize: "1rem", fontWeight: 600, margin: "0 0 1rem" },
  summaryGrid: { display: "flex", gap: "1rem", flexWrap: "wrap" },
  summaryItem: { background: "#f9fafb", borderRadius: 6, padding: "0.75rem 1rem", minWidth: 120 },
  summaryCategory: { fontSize: "0.8rem", color: "#6b7280", marginBottom: 4 },
  summaryTotal: { fontSize: "1.1rem", fontWeight: 700 },
  form: { display: "flex", gap: "0.5rem", flexWrap: "wrap", alignItems: "center" },
  input: { padding: "0.5rem 0.75rem", border: "1px solid #d1d5db", borderRadius: 6, fontSize: "0.9rem", flex: 1, minWidth: 120 },
  button: { padding: "0.5rem 1.25rem", background: "#111", color: "#fff", border: "none", borderRadius: 6, cursor: "pointer", fontWeight: 600 },
  error: { color: "#dc2626", marginTop: "0.5rem", fontSize: "0.875rem" },
  listHeader: { display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "1rem" },
  empty: { color: "#9ca3af", fontSize: "0.9rem" },
  table: { width: "100%", borderCollapse: "collapse" },
  th: { textAlign: "left", padding: "0.5rem 0.75rem", fontSize: "0.8rem", color: "#6b7280", borderBottom: "1px solid #e5e7eb" },
  tr: { borderBottom: "1px solid #f3f4f6" },
  td: { padding: "0.6rem 0.75rem", fontSize: "0.9rem" },
  deleteBtn: { background: "none", border: "none", cursor: "pointer", color: "#9ca3af", fontSize: "1rem", padding: "0 0.25rem" },
};
