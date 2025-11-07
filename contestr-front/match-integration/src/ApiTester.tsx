// ApiTester.tsx
import { useState } from 'react';
import { getCoreSportList } from '../client/sdk.gen';

export const ApiTester = () => {
    const [result, setResult] = useState<any>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // @ts-ignore
    const testConnection = async () => {
        setLoading(true);
        setError(null);
        try {
            console.log('🧪 Тестирование подключения к API...');
            const response = await getCoreSportList({});
            setResult(response);
            console.log('✅ API тест успешен:', response);
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Unknown error';
            setError(errorMessage);
            console.error('❌ API тест ошибка:', err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{
            padding: '20px',
            border: '1px solid #ccc',
            margin: '10px',
            borderRadius: '8px',
            background: '#f9f9f9'
        }}>
            <h3>🧪 API Connection Tester</h3>
            <button
                onClick={testConnection}
                disabled={loading}
                style={{
                    padding: '10px 15px',
                    background: loading ? '#ccc' : '#6466fd',
                    color: 'white',
                    border: 'none',
                    borderRadius: '5px',
                    cursor: loading ? 'not-allowed' : 'pointer'
                }}
            >
                {loading ? 'Testing...' : 'Test API Connection'}
            </button>

            {error && (
                <div style={{
                    color: '#d32f2f',
                    marginTop: '10px',
                    padding: '10px',
                    background: '#ffebee',
                    borderRadius: '4px'
                }}>
                    <strong>Error:</strong> {error}
                </div>
            )}

            {result && (
                <div style={{
                    marginTop: '10px',
                    padding: '10px',
                    background: '#e8f5e9',
                    borderRadius: '4px'
                }}>
                    <strong>Success! Response:</strong>
                    <pre style={{
                        background: '#fff',
                        padding: '10px',
                        borderRadius: '4px',
                        overflow: 'auto',
                        fontSize: '12px'
                    }}>
            {JSON.stringify(result, null, 2)}
          </pre>
                </div>
            )}
        </div>
    );
};