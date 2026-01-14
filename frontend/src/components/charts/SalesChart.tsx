import { Line } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend } from 'chart.js';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

export const SalesChart = ({ data, period }: { data: any[], period: string }) => {
  if (!data || data.length === 0) return <div style={{textAlign: 'center', padding: '2rem', color: 'var(--text-light)'}}>Aucune donnée de vente pour cette période.</div>;

  const chartData = {
    labels: data.map(d => d.label),
    datasets: [
      {
        label: 'Chiffre d\'affaires (CFA)',
        data: data.map(d => d.revenue),
        borderColor: 'rgb(53, 162, 235)',
        backgroundColor: 'rgba(53, 162, 235, 0.5)',
        yAxisID: 'y',
      },
      {
        label: 'Quantité Vendue',
        data: data.map(d => d.quantity),
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.5)',
        yAxisID: 'y1',
      },
    ],
  };

  const options = {
    responsive: true,
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
    stacked: false,
    plugins: {
      title: {
        display: true,
        text: `Évolution des Ventes (${period === 'day' ? 'Aujourd\'hui' : period === 'week' ? '7 Derniers Jours' : period === 'month' ? 'Ce Mois' : 'Cette Année'})`,
      },
    },
    scales: {
      y: {
        type: 'linear' as const,
        display: true,
        position: 'left' as const,
        title: { display: true, text: 'Revenu (CFA)' }
      },
      y1: {
        type: 'linear' as const,
        display: true,
        position: 'right' as const,
        grid: {
          drawOnChartArea: false,
        },
        title: { display: true, text: 'Quantité' }
      },
    },
  };

  return <Line options={options} data={chartData} />;
};
