import React, { useState } from 'react';
import {
  Download,
  Upload,
  Copy,
  Check,
  FileSpreadsheet,
  Database,
  Smartphone,
  Share2,
  Trash2,
  Info
} from 'lucide-react';
import {
  getSyncStats,
  exportDatabaseBinary,
  importDatabaseBinary,
  exportSyncPackage,
  importAndMergeSyncPackage,
  exportCSV,
  clearAllData
} from '../db/sqlite';
import { SyncStats } from '../types';

interface SettingsViewProps {
  onDataChanged: () => void;
}

export const SettingsView: React.FC<SettingsViewProps> = ({ onDataChanged }) => {
  const [stats, setStats] = useState<SyncStats>(() => getSyncStats());
  const [copied, setCopied] = useState<boolean>(false);
  const [syncStatusMsg, setSyncStatusMsg] = useState<string>('');
  const [isProcessing, setIsProcessing] = useState<boolean>(false);

  const refreshStats = () => {
    setStats(getSyncStats());
  };

  const handleCopyDeviceId = () => {
    navigator.clipboard.writeText(stats.deviceId);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleDownloadSQLite = () => {
    try {
      const binary = exportDatabaseBinary();
      const blob = new Blob([binary], { type: 'application/x-sqlite3' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `wallet_${new Date().toISOString().slice(0, 10)}.sqlite`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      alert(`Export error: ${err}`);
    }
  };

  const handleUploadSQLite = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!confirm('This will replace the current SQLite database with the selected file. Continue?')) {
      e.target.value = '';
      return;
    }

    setIsProcessing(true);
    const reader = new FileReader();
    reader.onload = async () => {
      try {
        const buffer = new Uint8Array(reader.result as ArrayBuffer);
        await importDatabaseBinary(buffer);
        refreshStats();
        onDataChanged();
        setSyncStatusMsg('SQLite database restored successfully!');
      } catch (err) {
        alert(`Failed to restore SQLite database: ${err}`);
      } finally {
        setIsProcessing(false);
      }
    };
    reader.readAsArrayBuffer(file);
    e.target.value = '';
  };

  const handleExportSyncJSON = () => {
    try {
      const pkg = exportSyncPackage();
      const jsonStr = JSON.stringify(pkg, null, 2);
      const blob = new Blob([jsonStr], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `wallet_sync_${new Date().toISOString().slice(0, 10)}.json`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      alert(`Sync export error: ${err}`);
    }
  };

  const handleImportSyncJSON = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsProcessing(true);
    const reader = new FileReader();
    reader.onload = async () => {
      try {
        const pkg = JSON.parse(reader.result as string);
        if (!pkg || !Array.isArray(pkg.records)) {
          throw new Error('Invalid sync package format');
        }
        const result = importAndMergeSyncPackage(pkg);
        refreshStats();
        onDataChanged();
        setSyncStatusMsg(
          `Merged cleanly: ${result.inserted} added, ${result.updated} updated, ${result.skipped} already current.`
        );
      } catch (err) {
        alert(`Failed to import sync JSON: ${err}`);
      } finally {
        setIsProcessing(false);
      }
    };
    reader.readAsText(file);
    e.target.value = '';
  };

  const handleExportCSV = () => {
    try {
      const csvStr = exportCSV();
      const blob = new Blob([csvStr], { type: 'text/csv;charset=utf-8;' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `wallet_expenses_${new Date().toISOString().slice(0, 10)}.csv`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      alert(`CSV export error: ${err}`);
    }
  };

  return (
    <div className="space-y-4">
      {/* Device & Sync Header */}
      <div className="bg-surface rounded-3xl p-5 border border-surface-border">
        <div className="flex items-center gap-2 mb-3">
          <Smartphone className="w-4 h-4 text-blue-400" />
          <h3 className="text-sm font-bold text-white">Device & Sync Identity</h3>
        </div>

        <p className="text-xs text-zinc-400 mb-3">
          Each device has a unique ID used for distributed Last-Write-Wins (LWW) conflict resolution when merging phone and desktop records.
        </p>

        <div className="bg-surface-elevated rounded-2xl p-3 border border-surface-border flex items-center justify-between">
          <div className="min-w-0 flex-1 pr-2">
            <div className="text-[10px] text-zinc-500 uppercase font-semibold">Device UUID</div>
            <div className="text-xs font-mono text-zinc-300 truncate">{stats.deviceId}</div>
          </div>
          <button
            onClick={handleCopyDeviceId}
            className="p-2 rounded-xl bg-surface border border-surface-border text-zinc-300 hover:text-white active:scale-90 transition-all flex items-center gap-1 text-xs"
          >
            {copied ? <Check className="w-4 h-4 text-emerald-400" /> : <Copy className="w-4 h-4" />}
          </button>
        </div>

        {/* Database Stats */}
        <div className="grid grid-cols-3 gap-2 mt-3 pt-3 border-t border-surface-border/60 text-center">
          <div className="p-2 rounded-xl bg-surface-elevated">
            <div className="text-[10px] text-zinc-400">Active Records</div>
            <div className="text-sm font-bold text-white mt-0.5">{stats.activeRecords}</div>
          </div>
          <div className="p-2 rounded-xl bg-surface-elevated">
            <div className="text-[10px] text-zinc-400">Tombstones</div>
            <div className="text-sm font-bold text-zinc-400 mt-0.5">{stats.deletedRecords}</div>
          </div>
          <div className="p-2 rounded-xl bg-surface-elevated">
            <div className="text-[10px] text-zinc-400">Total in DB</div>
            <div className="text-sm font-bold text-blue-400 mt-0.5">{stats.totalRecords}</div>
          </div>
        </div>
      </div>

      {/* Sync Status Banner */}
      {syncStatusMsg && (
        <div className="p-3 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 text-xs text-emerald-300 flex items-center gap-2">
          <Check className="w-4 h-4 flex-shrink-0" />
          <span>{syncStatusMsg}</span>
        </div>
      )}

      {/* Cross-Device Sync JSON Package */}
      <div className="bg-surface rounded-3xl p-5 border border-surface-border space-y-3">
        <div className="flex items-center gap-2">
          <Share2 className="w-4 h-4 text-indigo-400" />
          <h3 className="text-sm font-bold text-white">Cross-Device Data Sync</h3>
        </div>
        <p className="text-xs text-zinc-400">
          Transfer records between your phone and desktop. Merging uses timestamps to keep the newer version and preserves both sets without duplicates.
        </p>

        <div className="grid grid-cols-2 gap-2 pt-1">
          <button
            onClick={handleExportSyncJSON}
            disabled={isProcessing}
            className="p-3 rounded-2xl bg-surface-elevated border border-surface-border text-xs font-semibold text-zinc-200 hover:text-white hover:bg-zinc-800 active:scale-95 transition-all flex flex-col items-center gap-1.5"
          >
            <Download className="w-4 h-4 text-indigo-400" />
            Export Sync JSON
          </button>

          <label className="p-3 rounded-2xl bg-surface-elevated border border-surface-border text-xs font-semibold text-zinc-200 hover:text-white hover:bg-zinc-800 active:scale-95 transition-all flex flex-col items-center gap-1.5 cursor-pointer text-center">
            <Upload className="w-4 h-4 text-indigo-400" />
            Import & Merge
            <input
              type="file"
              accept=".json"
              onChange={handleImportSyncJSON}
              className="hidden"
            />
          </label>
        </div>
      </div>

      {/* Direct SQLite Database Binary File Operations */}
      <div className="bg-surface rounded-3xl p-5 border border-surface-border space-y-3">
        <div className="flex items-center gap-2">
          <Database className="w-4 h-4 text-emerald-400" />
          <h3 className="text-sm font-bold text-white">Native SQLite Database</h3>
        </div>
        <p className="text-xs text-zinc-400">
          Download the actual binary <code className="text-emerald-400">.sqlite</code> database file or restore a previous database file into your browser.
        </p>

        <div className="grid grid-cols-2 gap-2 pt-1">
          <button
            onClick={handleDownloadSQLite}
            disabled={isProcessing}
            className="p-3 rounded-2xl bg-surface-elevated border border-surface-border text-xs font-semibold text-zinc-200 hover:text-white hover:bg-zinc-800 active:scale-95 transition-all flex flex-col items-center gap-1.5"
          >
            <Download className="w-4 h-4 text-emerald-400" />
            Download .sqlite
          </button>

          <label className="p-3 rounded-2xl bg-surface-elevated border border-surface-border text-xs font-semibold text-zinc-200 hover:text-white hover:bg-zinc-800 active:scale-95 transition-all flex flex-col items-center gap-1.5 cursor-pointer text-center">
            <Upload className="w-4 h-4 text-emerald-400" />
            Restore .sqlite
            <input
              type="file"
              accept=".sqlite,.db,.sqlite3"
              onChange={handleUploadSQLite}
              className="hidden"
            />
          </label>
        </div>
      </div>

      {/* CSV Export for Spreadsheets */}
      <div className="bg-surface rounded-3xl p-5 border border-surface-border">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <FileSpreadsheet className="w-4 h-4 text-amber-400" />
            <div>
              <h3 className="text-sm font-bold text-white">Spreadsheet CSV</h3>
              <p className="text-[11px] text-zinc-400">Export table to open in Excel / Google Sheets</p>
            </div>
          </div>
          <button
            onClick={handleExportCSV}
            className="px-3 py-2 rounded-xl bg-surface-elevated border border-surface-border text-xs font-semibold text-zinc-200 hover:text-white active:scale-95 transition-all flex items-center gap-1.5"
          >
            <Download className="w-3.5 h-3.5 text-amber-400" />
            Export CSV
          </button>
        </div>
      </div>

      {/* Danger Zone: Reset Data */}
      <div className="bg-surface rounded-3xl p-5 border border-rose-500/20">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Trash2 className="w-4 h-4 text-rose-400" />
            <div>
              <h3 className="text-sm font-bold text-white">Reset Local Database</h3>
              <p className="text-[11px] text-zinc-400">Wipe all stored transactions on this device</p>
            </div>
          </div>
          <button
            onClick={() => {
              if (confirm('Are you sure you want to delete all transactions from this device? This cannot be undone.')) {
                clearAllData();
                refreshStats();
                onDataChanged();
                setSyncStatusMsg('All local transactions cleared.');
              }
            }}
            className="px-3 py-2 rounded-xl bg-rose-600/15 border border-rose-500/30 text-xs font-semibold text-rose-300 hover:bg-rose-600 hover:text-white active:scale-95 transition-all"
          >
            Clear Data
          </button>
        </div>
      </div>

      {/* GitHub Pages & Local Storage Note */}
      <div className="p-4 rounded-2xl bg-zinc-900/60 border border-surface-border text-zinc-500 text-xs flex gap-2.5">
        <Info className="w-4 h-4 flex-shrink-0 text-zinc-400 mt-0.5" />
        <p>
          Data is stored locally on this device in SQLite and IndexedDB. It works completely offline and will never be uploaded to any third-party cloud.
        </p>
      </div>
    </div>
  );
};
